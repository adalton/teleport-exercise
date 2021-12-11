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

package jobmanager

import (
	"os/exec"
	"sync"

	"github.com/adalton/teleport-exercise/pkg/cgroup/v1"
	"github.com/adalton/teleport-exercise/pkg/io"
	"github.com/google/uuid"
)

type job struct {
	mutex        sync.Mutex
	id           uuid.UUID
	name         string
	cgroupSet    *cgroup.Set
	cmd          exec.Cmd
	stdoutBuffer io.OutputBuffer
	stderrBuffer io.OutputBuffer
}

func NewJob(name string, cgroupSet *cgroup.Set, command string, args ...string) *job {
	return NewJobDetailed(name, cgroupSet, io.NewMemoryBuffer(), io.NewMemoryBuffer(), command, args...)
}

func NewJobDetailed(
	name string,
	cgroupSet *cgroup.Set,
	stdoutBuffer io.OutputBuffer,
	stderrBuffer io.OutputBuffer,
	command string,
	args ...string,
) *job {

	return &job{
		id:           uuid.New(),
		name:         name,
		cgroupSet:    cgroupSet,
		cmd:          *exec.Command(""),
		stdoutBuffer: stdoutBuffer,
		stderrBuffer: stderrBuffer,
	}
}
