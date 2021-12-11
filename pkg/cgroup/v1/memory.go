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

package cgroup

import "github.com/adalton/teleport-exercise/pkg/adaptation/os"

const (
	MemoryLimitInBytesFilename = "memory.limit_in_bytes"
)

// memory configures the memory cgroup controller.
type memory struct {
	base
	limit *string
}

func NewMemoryController() *memory {
	return NewMemoryControllerDetailed(nil)
}

func NewMemoryControllerDetailed(osAdapter *os.Adapter) *memory {
	return &memory{
		base: newBase("memory", osAdapter),
	}
}

func (m *memory) SetLimit(value string) *memory {
	m.limit = &value

	return m
}

func (m *memory) Apply(path string) error {
	if m.limit != nil {
		if err := m.write([]byte(*m.limit), "%s/%s", path, MemoryLimitInBytesFilename); err != nil {
			return err
		}
	}

	return nil
}
