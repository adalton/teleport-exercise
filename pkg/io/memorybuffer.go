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
	"fmt"
	"sync"
)

// The default initial capacity of a memory buffer created using NewMemoryBuffer.
const DefaultInitialMemoryBufferCapacity = 4096

// MemoryBuffer is an in-memory buffer of bytes.  This keeps the data even after
// it has been read to enable multiple clients to read what is written to this
// buffer.
type MemoryBuffer struct {
	mutex   sync.RWMutex
	content []byte
	closed  bool
}

// NewMemoryBuffer creates and returns a MemoryBuffer with an initial capacity
// of DefaultMemoryBufferCapacity.
func NewMemoryBuffer() *MemoryBuffer {
	return NewMemoryBufferDetailed(DefaultInitialMemoryBufferCapacity)
}

// NewMemoryBufferDetailed creates and returns a MemoryBuffer with the given
// initialCapacity.
func NewMemoryBufferDetailed(initialCapacity int) *MemoryBuffer {
	return &MemoryBuffer{
		content: make([]byte, 0, initialCapacity),
	}
}

// Write appends newContent to this MemoryBuffer.  The returned bytesWritten is
// always len(newContent).  Compatible with io.Writer.
func (b *MemoryBuffer) Write(newContent []byte) (bytesWritten int, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.closed {
		return 0, fmt.Errorf("cannot write to a closed MemoryBuffer")
	}

	b.content = append(b.content, newContent...)

	return len(newContent), nil
}

// ReadAt reads len(p) bytes into outputBuffer starting at the given offset in
// the underlying buffer.  It returns the number of bytes read
// (0 <= bytesRead <= len(outputBuffer)).  The returned error is always nil.
// Compatible with io.ReadAt.
func (b *MemoryBuffer) ReadAt(outputBuffer []byte, offset int64) (bytesRead int, err error) {
	if len(outputBuffer) == 0 {
		return 0, nil
	}

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return copy(outputBuffer, b.content[offset:]), nil
}

// Close closes this MemoryBuffer.  Once the MemoryBuffer is closed, it will
// accept no additional writes.  The returned error is always nil.
func (b *MemoryBuffer) Close() error {

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.closed = true

	return nil
}

// Size returns the current size of this MemoryBuffer.
func (b *MemoryBuffer) Size() int64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return int64(len(b.content))
}

// Closed returns true if this MemoryBuffer has been closed, false otherwise.
func (b *MemoryBuffer) Closed() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.closed
}
