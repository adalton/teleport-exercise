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
)

// DefaultMaxBufferSize is the default maximum number of bytes that a ByteStream
// will read from the underlying buffer at at once.
const DefaultMaxBufferSize = 1024

// ByteStream enables clients to stream bytes from an OutputBuffer.
type ByteStream struct {
	buffer      OutputBuffer
	channel     chan []byte
	maxReadSize int

	// mutex guards the fields that follow
	mutex            sync.Mutex
	goroutineStarted bool
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

	bailEarly := true

	func() {
		b.mutex.Lock()
		defer b.mutex.Unlock()

		if !b.goroutineStarted {
			bailEarly = false
			b.goroutineStarted = true
		}
	}()

	// If some other call to Stream has already created the goroutine, then
	// there's nothing for this call to do other than to return the channel.
	if bailEarly {
		return b.channel
	}

	go func() {
		var nextByte int64
		readBuffer := make([]byte, b.maxReadSize)

		for {
			bufferSize, _ := b.buffer.waitForChange(nextByte)

			// At this point the underlying buffer could be closed, there could
			// be new bytes in the buffer to process, or both.

			if bufferSize == nextByte {
				// No new bytes to process; the buffer must be closed and this
				// streamer must have consumed all the bytes that were written
				// to the buffer. Terminate the goroutine.
				close(b.channel)
				return
			}

			if n, err := b.buffer.ReadAt(readBuffer, nextByte); err != nil {
				// If ReadAt fails, we'l assume the buffer is in a bad state
				// and that future reads would also fail.

				log.Printf("Unexpected failure reading from underlying buffer: %v", err)

				close(b.channel)
				return
			} else if n > 0 {
				// Craete a copy here because we're reusing readBuffer here.
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

	return b.channel
}
