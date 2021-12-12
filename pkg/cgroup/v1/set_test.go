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
	"github.com/adalton/teleport-exercise/pkg/cgroup/v1/cgrouptest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Set_Create_Success(t *testing.T) {
	jobId, _ := uuid.Parse("0b5183b8-b572-49c7-90c4-fffc775b7d7b")
	mkdirAllRecorder := ostest.MkdirAllMock{}
	removeRecorder := ostest.RemoveMock{}

	adapter := &os.Adapter{
		MkdirAllFn: mkdirAllRecorder.MkdirAll,
		RemoveFn:   removeRecorder.Remove,
	}

	controller := &cgrouptest.DummyController{ControllerName: "nil"}
	set := cgroup.NewSetDetailed(adapter, cgroup.DefaultBasePath, jobId, controller)

	err := set.Create()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(mkdirAllRecorder.Events))
	assert.Equal(t, 0, len(removeRecorder.Events))
	assert.Equal(t,
		fmt.Sprintf("%s/%s/jobs/%s",
			cgroup.DefaultBasePath,
			controller.Name(),
			jobId.String(),
		),
		mkdirAllRecorder.Events[0].Path)
}

func Test_Set_Create_Failure(t *testing.T) {
	jobId, _ := uuid.Parse("0b5183b8-b572-49c7-90c4-fffc775b7d7b")
	mkdirAllRecorder := ostest.MkdirAllMock{}
	removeRecorder := ostest.RemoveMock{}

	adapter := &os.Adapter{
		MkdirAllFn: mkdirAllRecorder.MkdirAll,
		RemoveFn:   removeRecorder.Remove,
	}

	expectedError := fmt.Errorf("injected error")
	controller := &cgrouptest.DummyController{
		ControllerName:   "nil",
		ApplyReturnValue: expectedError,
	}
	set := cgroup.NewSetDetailed(adapter, cgroup.DefaultBasePath, jobId, controller)

	err := set.Create()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, len(removeRecorder.Events))
	assert.Equal(t,
		fmt.Sprintf("%s/%s/jobs/%s",
			cgroup.DefaultBasePath,
			controller.Name(),
			jobId.String(),
		),
		removeRecorder.Events[0].Path)
}

func Test_Set_Destroy_Success(t *testing.T) {
	jobId, _ := uuid.Parse("0b5183b8-b572-49c7-90c4-fffc775b7d7b")
	removeRecorder := ostest.RemoveMock{}

	adapter := &os.Adapter{
		RemoveFn: removeRecorder.Remove,
	}

	controller := &cgrouptest.DummyController{ControllerName: "nil"}
	set := cgroup.NewSetDetailed(adapter, cgroup.DefaultBasePath, jobId, controller)

	err := set.Destroy()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(removeRecorder.Events))
	assert.Equal(t,
		fmt.Sprintf("%s/%s/jobs/%s",
			cgroup.DefaultBasePath,
			controller.Name(),
			jobId.String(),
		),
		removeRecorder.Events[0].Path)
}

func Test_Set_Destroy_Failure(t *testing.T) {
	jobId, _ := uuid.Parse("0b5183b8-b572-49c7-90c4-fffc775b7d7b")
	injectedError := fmt.Errorf("injected error")
	removeRecorder := ostest.RemoveMock{
		NextError: injectedError,
	}

	adapter := &os.Adapter{
		RemoveFn: removeRecorder.Remove,
	}

	controller := &cgrouptest.DummyController{ControllerName: "nil"}
	set := cgroup.NewSetDetailed(adapter, cgroup.DefaultBasePath, jobId, controller)

	err := set.Destroy()

	assert.Error(t, err)
	assert.Equal(t, 1, len(removeRecorder.Events))
}

func Test_Set_TaskFiles(t *testing.T) {
	jobId, _ := uuid.Parse("0b5183b8-b572-49c7-90c4-fffc775b7d7b")

	controller := &cgrouptest.DummyController{ControllerName: "nil"}
	set := cgroup.NewSet(jobId, controller)

	taskFiles := set.TaskFiles()

	assert.Equal(t, 1, len(taskFiles))
	assert.Equal(t,
		fmt.Sprintf("%s/%s/jobs/%s/tasks",
			cgroup.DefaultBasePath,
			controller.Name(),
			jobId.String(),
		),
		taskFiles[0])
}
