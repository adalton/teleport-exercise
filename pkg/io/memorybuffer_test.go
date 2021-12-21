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

package io_test

import (
	"testing"

	"github.com/adalton/teleport-exercise/pkg/io"

	"github.com/stretchr/testify/assert"
)

func Test_MemoryBuffer_InitialSizeZero(t *testing.T) {
	b := io.NewMemoryBuffer()

	assert.Equal(t, int64(0), b.Size())
}

func Test_MemoryBuffer_ReadAt_FromEmptyBuffer(t *testing.T) {
	b := io.NewMemoryBuffer()
	result := make([]byte, 64)

	bytesRead, err := b.ReadAt(result, 0)

	assert.Nil(t, err)
	assert.Equal(t, 0, bytesRead)
}

func Test_MemoryBuffer_Write_EmptyBuffer(t *testing.T) {
	b := io.NewMemoryBuffer()

	bytesWritten, err := b.Write([]byte(""))

	assert.Nil(t, err)
	assert.Equal(t, 0, bytesWritten)
}

func Test_MemoryBuffer_WriteAfterClose(t *testing.T) {
	b := io.NewMemoryBuffer()

	err := b.Close()
	assert.Nil(t, err)
	assert.True(t, b.Closed())

	_, err = b.Write([]byte(""))

	assert.Error(t, err)
}

func Test_MemoryBuffer_ReadAt_FromBeginning(t *testing.T) {
	b := io.NewMemoryBuffer()
	result := make([]byte, 5)

	b.Write([]byte("abcdefghijklmnopqrstuvwxyz"))

	bytesRead, err := b.ReadAt(result, 0)

	assert.Nil(t, err)
	assert.Equal(t, len(result), bytesRead)
	assert.Equal(t, []byte("abcde"), result)
}

func Test_MemoryBuffer_ReadAt_FromMiddle(t *testing.T) {
	b := io.NewMemoryBuffer()
	result := make([]byte, 7)

	b.Write([]byte("abcdefghijklmnopqrstuvwxyz"))

	bytesRead, err := b.ReadAt(result, 2)

	assert.Nil(t, err)
	assert.Equal(t, len(result), bytesRead)
	assert.Equal(t, []byte("cdefghi"), result)
}

func Test_MemoryBuffer_ReadAt_FromEnd(t *testing.T) {
	b := io.NewMemoryBuffer()
	result := make([]byte, 7)

	b.Write([]byte("abcdefghijklmnopqrstuvwxyz"))

	bytesRead, err := b.ReadAt(result, 26-3)

	assert.Nil(t, err)
	assert.Equal(t, 3, bytesRead)
	assert.Equal(t, []byte("xyz"), result[0:bytesRead])
}

func Test_MemoryBuffer_ReadAt_All(t *testing.T) {
	b := io.NewMemoryBuffer()
	content := []byte("abcdefghijklmnopqrstuvwxyz")
	result := make([]byte, 26)

	b.Write(content)

	bytesRead, err := b.ReadAt(result, 0)

	assert.Nil(t, err)
	assert.Equal(t, len(content), bytesRead)
	assert.Equal(t, content, result)
}

func Test_MemoryBuffer_ReadAt_EmptyBuffer(t *testing.T) {
	b := io.NewMemoryBuffer()
	content := []byte("abcdefghijklmnopqrstuvwxyz")
	result := make([]byte, 0)

	b.Write(content)

	bytesRead, err := b.ReadAt(result, 0)

	assert.Nil(t, err)
	assert.Equal(t, 0, bytesRead)
	assert.Equal(t, 0, len(result))
}

func Test_MemoryBuffer_ReadAt_OffsetEqualToSize(t *testing.T) {
	b := io.NewMemoryBuffer()
	content := []byte("abcdefghijklmnopqrstuvwxyz")
	result := make([]byte, 5)

	b.Write(content)

	bytesRead, err := b.ReadAt(result, int64(len(content)))

	assert.Nil(t, err)
	assert.Equal(t, 0, bytesRead)
}

func Test_MemoryBuffer_ReadAt_OffsetGreaterThanSize(t *testing.T) {
	b := io.NewMemoryBuffer()
	content := []byte("abcdefghijklmnopqrstuvwxyz")
	result := make([]byte, 5)

	b.Write(content)

	_, err := b.ReadAt(result, int64(len(content)+1))

	assert.Error(t, err)
}
