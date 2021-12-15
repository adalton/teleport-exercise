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

package cgroupv1

import (
	"fmt"

	"github.com/adalton/teleport-exercise/pkg/adaptation/os"
)

const (
	MemoryLimitInBytesFilename = "memory.limit_in_bytes"
)

// MemoryController configures the MemoryController cgroup controller.
type MemoryController struct {
	OsAdapter *os.Adapter
	Limit     string
}

func (MemoryController) Name() string {
	return "memory"
}

func (m *MemoryController) Apply(path string) error {
	if m.Limit != "" {
		filename := fmt.Sprintf("%s/%s", path, MemoryLimitInBytesFilename)
		if err := m.OsAdapter.WriteFile(filename, []byte(m.Limit), os.FileMode(0644)); err != nil {
			return err
		}
	}

	return nil
}
