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
	mkdirAllRecorder := ostest.MkdirAllRecorder{}
	removeRecorder := ostest.RemoveRecorder{}

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
	mkdirAllRecorder := ostest.MkdirAllRecorder{}
	removeRecorder := ostest.RemoveRecorder{}

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
	removeRecorder := ostest.RemoveRecorder{}

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
	removeRecorder := ostest.RemoveRecorder{
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
