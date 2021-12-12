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

package command_test

import (
	"fmt"
	"testing"

	"github.com/adalton/teleport-exercise/pkg/adaptation/os"
	"github.com/adalton/teleport-exercise/pkg/adaptation/os/ostest"
	"github.com/adalton/teleport-exercise/pkg/adaptation/syscall"
	"github.com/adalton/teleport-exercise/pkg/adaptation/syscall/syscalltest"
	"github.com/adalton/teleport-exercise/pkg/command"

	"github.com/stretchr/testify/assert"
)

func Test_Cgexec_WriteCgroupFiles_Success(t *testing.T) {
	writeFileRecorder := &ostest.WriteFileMock{}
	pidGenerator := &ostest.GetpidMock{
		Pid: 1234,
	}
	osa := &os.Adapter{
		WriteFileFn: writeFileRecorder.WriteFile,
		GetpidFn:    pidGenerator.Getpid,
	}

	sc := &syscall.Adapter{
		ExecFn: (&syscalltest.ExecMock{}).Exec,
	}

	cgfile := "/sys/fs/cgroup/cpu/job/1e71d42d-b7e2-4f1c-893f-b16415b96e1a/tasks"

	args := []string{
		"nameOfTheTool",
		cgfile,
		"--",
		"ignored in this test",
	}

	_ = command.CgexecDetailed(args, osa, sc)

	assert.Equal(t, 1, len(writeFileRecorder.Events))
	assert.Equal(t, cgfile, writeFileRecorder.Events[0].Name)
	assert.Equal(t, fmt.Sprintf("%d", pidGenerator.Pid), string(writeFileRecorder.Events[0].Data))
}

func Test_Cgexec_WriteCgroupFiles_Failure(t *testing.T) {
	expectedError := fmt.Errorf("injected error")
	writeFileRecorder := &ostest.WriteFileMock{
		NextError: expectedError,
	}
	osa := &os.Adapter{
		WriteFileFn: writeFileRecorder.WriteFile,
		GetpidFn:    (&ostest.GetpidMock{}).Getpid,
	}

	sc := &syscall.Adapter{
		ExecFn: (&syscalltest.ExecMock{}).Exec,
	}

	cgfile := "/sys/fs/cgroup/cpu/job/1e71d42d-b7e2-4f1c-893f-b16415b96e1a/tasks"

	args := []string{
		"nameOfTheTool",
		cgfile,
		"--",
		"ignored in this test",
	}

	err := command.CgexecDetailed(args, osa, sc)

	assert.Equal(t, expectedError, err)
}

func Test_Cgexec_Exec(t *testing.T) {
	env := []string{"x=y"}
	envGen := &ostest.EnvironMock{
		Environment: env,
	}

	osa := &os.Adapter{
		WriteFileFn: (&ostest.WriteFileMock{}).WriteFile,
		GetpidFn:    (&ostest.GetpidMock{}).Getpid,
		EnvironFn:   envGen.Environ,
	}

	execRecorder := &syscalltest.ExecMock{}
	sc := &syscall.Adapter{
		ExecFn: execRecorder.Exec,
	}

	commandName := "commandName"
	commandArgs := []string{"arg1", "arg2", "--", "arg3"}
	args := []string{
		"nameOfTheTool",
		"--",
		commandName,
	}
	args = append(args, commandArgs...)

	var argv []string
	argv = append(argv, commandName)
	argv = append(argv, commandArgs...)

	err := command.CgexecDetailed(args, osa, sc)

	assert.Error(t, err)
	assert.Equal(t, commandName, execRecorder.Argv0)
	assert.Equal(t, argv, execRecorder.Argv)
	assert.Equal(t, env, execRecorder.Envv)
}
