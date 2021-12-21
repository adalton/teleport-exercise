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

package serverv1_test

import (
	"context"
	"testing"

	"github.com/adalton/teleport-exercise/pkg/jobmanager"
	"github.com/adalton/teleport-exercise/pkg/jobmanager/jobmanagertest"
	"github.com/adalton/teleport-exercise/server/jobmanager/serverv1"
	"github.com/adalton/teleport-exercise/server/jobmanager/serverv1/testserverv1"
	"github.com/adalton/teleport-exercise/service/jobmanager/jobmanagerv1"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_jobmanagerServer_Start_NoUserID(t *testing.T) {
	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)

	_, err := server.Start(context.Background(), &jobmanagerv1.JobCreationRequest{})

	assert.Error(t, err)
	assert.ErrorIs(t, err, jobmanager.ErrUnauthenticated)
}

func Test_jobmanagerServer_Start_WithUserID(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)

	ctx := serverv1.AttachUserIDToContext(context.Background(), "user1")

	job, err := server.Start(ctx, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	_, parseErr := uuid.Parse(job.Id.Id)

	assert.Nil(t, parseErr)
	assert.Equal(t, jobName, job.Name)
}

func Test_jobmanagerServer_Start_NameExists(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)

	ctx := serverv1.AttachUserIDToContext(context.Background(), "user1")

	_, err := server.Start(ctx, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	_, err = server.Start(ctx, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})

	assert.ErrorIs(t, err, jobmanager.ErrJobExists)
}

func Test_jobmanagerServer_Stop_NoUserID(t *testing.T) {
	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)

	_, err := server.Stop(context.Background(), &jobmanagerv1.JobID{Id: "b13620d4-db7f-46d5-b445-b29af0f87d2c"})

	assert.ErrorIs(t, err, jobmanager.ErrUnauthenticated)
}

func Test_jobmanagerServer_Stop_MalformedJobID(t *testing.T) {
	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	ctx := serverv1.AttachUserIDToContext(context.Background(), "user1")

	_, err := server.Stop(ctx, &jobmanagerv1.JobID{Id: "not-a-valid-id"})

	assert.ErrorIs(t, err, jobmanager.ErrInvalidJobID)
}

func Test_jobmanagerServer_Stop_JobDoesNotExist(t *testing.T) {
	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	ctx := serverv1.AttachUserIDToContext(context.Background(), "user1")

	_, err := server.Stop(ctx, &jobmanagerv1.JobID{Id: "eeafbe44-348f-47ba-ba2b-3e013ee8bb85"})

	assert.ErrorIs(t, err, jobmanager.ErrJobNotFound)
}

func Test_jobmanagerServer_Stop_JobExists(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	ctx := serverv1.AttachUserIDToContext(context.Background(), "user1")

	job, err := server.Start(ctx, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	_, err = server.Stop(ctx, &jobmanagerv1.JobID{Id: job.Id.Id})
	assert.Nil(t, err)
}

func Test_jobmanagerServer_Query_NoUserID(t *testing.T) {
	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)

	_, err := server.Query(context.Background(), &jobmanagerv1.JobID{Id: "3e3d8936-5fd7-46bb-9fd2-8423c607a0b2"})

	assert.ErrorIs(t, err, jobmanager.ErrUnauthenticated)
}

func Test_jobmanagerServer_Query_MalformedJobID(t *testing.T) {
	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	ctx := serverv1.AttachUserIDToContext(context.Background(), "user1")

	_, err := server.Query(ctx, &jobmanagerv1.JobID{Id: "not-a-valid-jobID"})

	assert.ErrorIs(t, err, jobmanager.ErrInvalidJobID)
}

func Test_jobmanagerServer_Query_JobExists(t *testing.T) {
	const (
		owner       = "user1"
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	ctx := serverv1.AttachUserIDToContext(context.Background(), owner)

	job, err := server.Start(ctx, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	jobStatus, err := server.Query(ctx, &jobmanagerv1.JobID{Id: job.Id.Id})

	assert.Nil(t, err)
	assert.Equal(t, jobName, jobStatus.Job.Name)
	assert.Equal(t, job.Id, jobStatus.Job.Id)
	assert.Equal(t, owner, jobStatus.Owner)
	assert.True(t, jobStatus.IsRunning)
	assert.Equal(t, int32(jobmanagertest.DefaultPID), jobStatus.Pid)
	assert.Equal(t, int32(jobmanagertest.DefaultSignalWhileRunning), jobStatus.SignalNumber)
	assert.Equal(t, "", jobStatus.ErrorMessage)
}

func Test_jobmanagerServer_List_NoUserID(t *testing.T) {
	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)

	_, err := server.List(context.Background(), &jobmanagerv1.NilMessage{})

	assert.Error(t, err, jobmanager.ErrUnauthenticated)
}

func Test_jobmanagerServer_List_NoJobs(t *testing.T) {
	const owner = "user1"
	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	ctx := serverv1.AttachUserIDToContext(context.Background(), owner)

	jobList, err := server.List(ctx, &jobmanagerv1.NilMessage{})

	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobList.JobStatusList))
}

func Test_jobmanagerServer_List_JobExists(t *testing.T) {
	const (
		owner       = "user1"
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	ctx := serverv1.AttachUserIDToContext(context.Background(), owner)

	job, err := server.Start(ctx, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	jobList, err := server.List(ctx, &jobmanagerv1.NilMessage{})

	assert.Nil(t, err)
	assert.Equal(t, 1, len(jobList.JobStatusList))
	assert.Equal(t, jobName, jobList.JobStatusList[0].Job.Name)
	assert.Equal(t, job.Id, jobList.JobStatusList[0].Job.Id)
	assert.Equal(t, owner, jobList.JobStatusList[0].Owner)
	assert.True(t, jobList.JobStatusList[0].IsRunning)
	assert.Equal(t, int32(jobmanagertest.DefaultPID), jobList.JobStatusList[0].Pid)
	assert.Equal(t, int32(jobmanagertest.DefaultSignalWhileRunning), jobList.JobStatusList[0].SignalNumber)
	assert.Equal(t, "", jobList.JobStatusList[0].ErrorMessage)
}

func Test_jobmanagerServer_Stream_NoUserID(t *testing.T) {
	mockServer := &testserverv1.MockJobmanagerStreamServer{
		NextContext: context.Background(),
	}

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)

	err := server.StreamOutput(&jobmanagerv1.StreamOutputRequest{}, mockServer)

	assert.ErrorIs(t, err, jobmanager.ErrUnauthenticated)
}

func Test_jobmanagerServer_MalformedJobID(t *testing.T) {
	ctx := serverv1.AttachUserIDToContext(context.Background(), "user1")
	mockServer := &testserverv1.MockJobmanagerStreamServer{
		NextContext: ctx,
	}
	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	req := &jobmanagerv1.StreamOutputRequest{
		JobID:        &jobmanagerv1.JobID{Id: "not-a-valid-jobID"},
		OutputStream: jobmanagerv1.OutputStream_OutputStream_STDOUT,
	}

	err := server.StreamOutput(req, mockServer)
	assert.ErrorIs(t, err, jobmanager.ErrInvalidJobID)
}

func Test_jobmanagerServer_InvalidStreamType(t *testing.T) {
	ctx := serverv1.AttachUserIDToContext(context.Background(), "user1")
	mockServer := &testserverv1.MockJobmanagerStreamServer{
		NextContext: ctx,
	}
	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	req := &jobmanagerv1.StreamOutputRequest{
		JobID:        &jobmanagerv1.JobID{Id: "1294326a-816a-4f13-8a8b-ff92c7a78984"},
		OutputStream: jobmanagerv1.OutputStream(-1),
	}

	err := server.StreamOutput(req, mockServer)
	assert.ErrorIs(t, err, jobmanager.ErrInvalidArgument)
}

func Test_jobmanagerServer_Stream_ContextCanceled(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}
	ctx := serverv1.AttachUserIDToContext(context.Background(), "user1")
	ctx, cancel := context.WithCancel(ctx)

	mockServer := &testserverv1.MockJobmanagerStreamServer{
		NextContext: ctx,
	}

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	job, err := server.Start(ctx, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	req := &jobmanagerv1.StreamOutputRequest{
		JobID:        &jobmanagerv1.JobID{Id: job.Id.Id},
		OutputStream: jobmanagerv1.OutputStream_OutputStream_STDERR,
	}

	cancel() // Intentially calling this here, not deferring it
	err = server.StreamOutput(req, mockServer)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func Test_jobmanagerServer_Stream_ReadSuccessfully(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}
	ctx := serverv1.AttachUserIDToContext(context.Background(), "user1")
	mockServer := &testserverv1.MockJobmanagerStreamServer{
		NextContext: ctx,
	}

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	job, err := server.Start(ctx, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	req := &jobmanagerv1.StreamOutputRequest{
		JobID:        &jobmanagerv1.JobID{Id: job.Id.Id},
		OutputStream: jobmanagerv1.OutputStream_OutputStream_STDOUT,
	}

	_, err = server.Stop(ctx, &jobmanagerv1.JobID{Id: job.Id.Id})
	assert.Nil(t, err)

	err = server.StreamOutput(req, mockServer)

	assert.Nil(t, err)
	require.NotNil(t, mockServer.LastJobOutput)
	assert.Equal(t, []byte(jobmanagertest.DefaultStandardOutput), mockServer.LastJobOutput.Output)
}

func Test_jobmanagerServer_Multitenant(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}
	ctxUser1 := serverv1.AttachUserIDToContext(context.Background(), "user1")

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	_, err := server.Start(ctxUser1, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	// User2 cannot set User1's jobs
	ctxUser2 := serverv1.AttachUserIDToContext(context.Background(), "user2")
	jobList, err := server.List(ctxUser2, &jobmanagerv1.NilMessage{})

	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobList.JobStatusList))
}

func Test_jobmanagerServer_Multitenant_Stop(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}
	ctxUser1 := serverv1.AttachUserIDToContext(context.Background(), "user1")

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	job, err := server.Start(ctxUser1, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	// User2 cannot set User1's jobs
	ctxUser2 := serverv1.AttachUserIDToContext(context.Background(), "user2")
	_, err = server.Stop(ctxUser2, &jobmanagerv1.JobID{Id: job.Id.Id})

	assert.ErrorIs(t, err, jobmanager.ErrJobNotFound)
}

func Test_jobmanagerServer_AdministratorCanSeeAllJobs(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}
	ctxUser1 := serverv1.AttachUserIDToContext(context.Background(), "user1")

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	_, err := server.Start(ctxUser1, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	ctxUser2 := serverv1.AttachUserIDToContext(context.Background(), "user2")
	_, err = server.Start(ctxUser2, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	// Superuser can see user1's and user2's jobs
	ctxAdmin := serverv1.AttachUserIDToContext(context.Background(), jobmanager.Superuser)
	jobList, err := server.List(ctxAdmin, &jobmanagerv1.NilMessage{})

	assert.Nil(t, err)
	assert.Equal(t, 2, len(jobList.JobStatusList))
	if jobList.JobStatusList[0].Owner == "user1" {
		assert.Equal(t, "user1", jobList.JobStatusList[0].Owner)
		assert.Equal(t, "user2", jobList.JobStatusList[1].Owner)
	} else {
		assert.Equal(t, "user2", jobList.JobStatusList[0].Owner)
		assert.Equal(t, "user1", jobList.JobStatusList[1].Owner)
	}
}

func Test_jobmanagerServer_User2CannotStopUser1sJob(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}
	ctxUser1 := serverv1.AttachUserIDToContext(context.Background(), "user1")

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	job, err := server.Start(ctxUser1, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	ctxUser2 := serverv1.AttachUserIDToContext(context.Background(), "user2")
	_, err = server.Stop(ctxUser2, &jobmanagerv1.JobID{Id: job.Id.Id})

	assert.ErrorIs(t, err, jobmanager.ErrJobNotFound)
}

func Test_jobmanagerServer_AdministratorCanStopUser1sJob(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}
	ctxUser1 := serverv1.AttachUserIDToContext(context.Background(), "user1")

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	job, err := server.Start(ctxUser1, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	ctxAdmin := serverv1.AttachUserIDToContext(context.Background(), jobmanager.Superuser)
	_, err = server.Stop(ctxAdmin, &jobmanagerv1.JobID{Id: job.Id.Id})
	assert.Nil(t, err)

	jobStatus, err := server.Query(ctxUser1, &jobmanagerv1.JobID{Id: job.Id.Id})
	assert.Nil(t, err)

	assert.False(t, jobStatus.IsRunning)
}

func Test_jobmanagerServer_User2CannotQueryUser1sJob(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}
	ctxUser1 := serverv1.AttachUserIDToContext(context.Background(), "user1")

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	job, err := server.Start(ctxUser1, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	ctxUser2 := serverv1.AttachUserIDToContext(context.Background(), "user2")
	_, err = server.Query(ctxUser2, &jobmanagerv1.JobID{Id: job.Id.Id})

	assert.ErrorIs(t, err, jobmanager.ErrJobNotFound)
}

func Test_jobmanagerServer_AdministratorCanQueryUser1sJob(t *testing.T) {
	const (
		jobName     = "myJob"
		programPath = "/bin/ls"
	)
	args := []string{"-l", "/"}
	ctxUser1 := serverv1.AttachUserIDToContext(context.Background(), "user1")

	jobManager := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	server := serverv1.NewJobManagerServerDetailed(jobManager)
	job, err := server.Start(ctxUser1, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   args,
	})
	assert.Nil(t, err)

	ctxAdmin := serverv1.AttachUserIDToContext(context.Background(), jobmanager.Superuser)
	jobStatus, err := server.Query(ctxAdmin, &jobmanagerv1.JobID{Id: job.Id.Id})

	assert.Nil(t, err)
	assert.Equal(t, job.Id.Id, jobStatus.Job.Id.Id)
}
