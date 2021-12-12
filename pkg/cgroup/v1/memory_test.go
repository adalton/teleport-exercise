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

package cgroup_test

import (
	"fmt"
	"testing"

	"github.com/adalton/teleport-exercise/pkg/adaptation/os"
	"github.com/adalton/teleport-exercise/pkg/adaptation/os/ostest"
	"github.com/adalton/teleport-exercise/pkg/cgroup/v1"
	"github.com/stretchr/testify/assert"
)

func Test_memory_Apply(t *testing.T) {
	path := "/sys/fs/cgroup/jobs/889f7cc2-9935-4773-aaa1-b94478abc923"
	writeRecorder := ostest.WriteFileMock{}
	adapter := &os.Adapter{
		WriteFileFn: writeRecorder.WriteFile,
	}

	limit := "500M"
	mem := cgroup.NewMemoryControllerDetailed(adapter).SetLimit(limit)
	mem.Apply(path)

	assert.Equal(t, 1, len(writeRecorder.Events))
	assert.Equal(t, fmt.Sprintf("%s/%s", path, cgroup.MemoryLimitInBytesFilename), writeRecorder.Events[0].Name)
	assert.Equal(t, []byte(limit), writeRecorder.Events[0].Data)
}
