/*
Copyright 2021 Andy Dalton
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package io

import (
	"log"
	"sync"
	"sync/atomic"
)

// DefaultMaxBufferSize is the default maximum number of bytes that a ByteStream
// will read from the underlying buffer at at once.
const DefaultMaxBufferSize = 1024

// ByteStream enables clients to stream bytes from an OutputBuffer.
type ByteStream struct {
	startGoroutine sync.Once
	buffer         OutputBuffer
	channel        chan []byte
	maxReadSize    int
	closed         int32 // 0 = not closed, 1 = closed
}

// NewByteStream creates and returns a new ByteStream associated with the
// given buffer with a read buffer size of DefaultMaxBufferSize.
func NewByteStream(buffer OutputBuffer) *ByteStream {
	return NewByteStreamDetailed(buffer, DefaultMaxBufferSize)
}

// NewByteStream creates and returns a new ByteStream associated with the
// given buffer with a read buffer size of the given maxReadSize.
func NewByteStreamDetailed(buffer OutputBuffer, maxReadSize int) *ByteStream {
	return &ByteStream{
		channel:     make(chan []byte),
		buffer:      buffer,
		maxReadSize: maxReadSize,
	}
}

// Stream returns a chanel that streams the content of the underlying OutputBuffer.
func (b *ByteStream) Stream() <-chan []byte {

	b.startGoroutine.Do(func() {
		go func() {
			defer close(b.channel)

			var nextByte int64
			readBuffer := make([]byte, b.maxReadSize)

			for {
				if b.Closed() {
					return
				}

				bufferSize, outputBufferClosed := b.buffer.waitForChange(nextByte)

				// At this point this ByteStream could be closed.
				// Also, the underlyingbuffer could be closed, there could be
				// new bytes in the buffer to process, or both.

				if b.Closed() || (bufferSize == nextByte && outputBufferClosed) {
					// No new bytes to process and the buffer is closed.  This
					// streamer must have consumed all the bytes that were written
					// to the buffer. Terminate the goroutine.
					return
				}

				if n, err := b.buffer.ReadAt(readBuffer, nextByte); err != nil {
					// If ReadAt fails, we'l assume the buffer is in a bad state
					// and that future reads would also fail.
					//
					// One way in which this could happen if nextByte is greater
					// than the size of the buffer.  That should never happen.
					// It might be worth panicing here --- at least during
					// internal testing --- to catch this error if it happens.

					log.Printf("Unexpected failure reading from underlying buffer: %v", err)
					return
				} else if n > 0 {
					// Create a copy here because we're reusing readBuffer here.
					bufToWrite := make([]byte, n)
					copy(bufToWrite, readBuffer[0:n])

					b.channel <- bufToWrite
					nextByte += int64(n)
				}

				// If n == 0, then it's possible that new bytes were written to the
				// buffer since we inspected its size above.  In that case, we'll
				// get the updated size and reevaluate on the next iteration.
			}
		}()

	})

	return b.channel
}

// Close marks this ByteStream as closed.  If there is a goroutine associated with
// this ByteStream, then it will close the stream and terminate before this
// function returns.
func (b *ByteStream) Close() {

	// Mark this ByteStream as closed if it's not already closed
	if !atomic.CompareAndSwapInt32(&b.closed, 0, 1) {
		// Already closed, nothing to do
		return
	}

	// If the goroutine hasn't yet been started, block it from being started and
	// maintain state to track whether it could still be running
	goroutineWasStarted := true
	b.startGoroutine.Do(func() { goroutineWasStarted = false })

	// If the goroutine wasn't started, then we still have an open channel to
	// which nothing has been written.  We can safely closed it here because
	// we've blocked the goroutine from starting, and that's the only other
	// thing that could close it.
	if !goroutineWasStarted {
		close(b.channel)
		return
	}

	// If goroutineWasStarted is true, then the goroutine associated with this
	// ByteStream was started at some point in time.  It may still be running,
	// or it may have already returned.  All exit paths from the goroutine close
	// the channel.
	//
	// This function updated b.closed to the closed state, so if the goroutine
	// is potentially still running then we need to give it the opportunity to
	// notice that this ByteStream is closed and return, thus closing the
	// channel.  However, the goroutine could be blocked on a write to the
	// channel.
	//
	// If we read from the channel and it's already closed, we can noticed and
	// return.  If we read from the channel and it returns data, then that'll
	// unblock the goroutine that was blocked on a write to the channel, enable
	// it to see that this ByteStream is closed, and close the channel.
	//
	// Once we know the channel is closed, then we can safely return.
	_, channelOpen := <-b.channel
	for channelOpen {
		_, channelOpen = <-b.channel
	}
}

// Closed returns true if this ByteStream is closed, false otherwise.
func (b *ByteStream) Closed() bool {
	return atomic.LoadInt32(&b.closed) > 0
}
