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

package serverv1

import (
	"context"

	"github.com/adalton/teleport-exercise/pkg/io"
	"github.com/adalton/teleport-exercise/pkg/jobmanager"
	"github.com/adalton/teleport-exercise/service/jobmanager/jobmanagerv1"
)

// jobmanagerServer implements the gRPC handler for the jobmanager service.
type jobmanagerServer struct {
	jobmanagerv1.UnimplementedJobManagerServer
	jm *jobmanager.Manager
}

// NewJobmanagerServer creates and returns a jobmanagerServer.
func NewJobmanagerServer() *jobmanagerServer {
	return NewJobManagerServerDetailed(jobmanager.NewManager())
}

// NewJobManagerServerDetailed creates a new jobmanagerServer with a custom
// underlying jobmanager.Manager.
func NewJobManagerServerDetailed(manager *jobmanager.Manager) *jobmanagerServer {
	return &jobmanagerServer{
		jm: manager,
	}
}

func (s *jobmanagerServer) Start(
	ctx context.Context,
	jcr *jobmanagerv1.JobCreationRequest,
) (*jobmanagerv1.Job, error) {

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	job, err := s.jm.Start(userID, jcr.GetName(), jcr.GetProgramPath(), jcr.GetArguments())
	if err != nil {
		return nil, err
	}

	jobResponse := &jobmanagerv1.Job{
		Id:   &jobmanagerv1.JobID{Id: job.ID().String()},
		Name: job.Name(),
	}

	return jobResponse, nil
}

func (s *jobmanagerServer) Stop(
	ctx context.Context,
	requestJobID *jobmanagerv1.JobID,
) (*jobmanagerv1.NilMessage, error) {

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = s.jm.Stop(userID, requestJobID.Id)
	if err != nil {
		return nil, err
	}

	return &jobmanagerv1.NilMessage{}, nil
}

func internalToExternalStatusV1(internalStatus *jobmanager.JobStatus) *jobmanagerv1.JobStatus {
	errMsg := ""

	if internalStatus.RunError != nil {
		errMsg = internalStatus.RunError.Error()
	}

	return &jobmanagerv1.JobStatus{
		Job: &jobmanagerv1.Job{
			Id: &jobmanagerv1.JobID{
				Id: internalStatus.ID,
			},
			Name: internalStatus.Name,
		},
		Owner:        internalStatus.Owner,
		IsRunning:    internalStatus.Running,
		Pid:          int32(internalStatus.Pid),
		ExitCode:     int32(internalStatus.ExitCode),
		SignalNumber: int32(internalStatus.SignalNum),
		ErrorMessage: errMsg,
	}
}

func (s *jobmanagerServer) Query(
	ctx context.Context,
	requestJobID *jobmanagerv1.JobID,
) (*jobmanagerv1.JobStatus, error) {

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	jobStatus, err := s.jm.Status(userID, requestJobID.Id)
	if err != nil {
		return nil, err
	}

	return internalToExternalStatusV1(jobStatus), nil
}

func (s *jobmanagerServer) List(
	ctx context.Context,
	_ *jobmanagerv1.NilMessage,
) (*jobmanagerv1.JobStatusList, error) {

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	statusList := s.jm.List(userID)

	responseStatusList := &jobmanagerv1.JobStatusList{
		JobStatusList: make([]*jobmanagerv1.JobStatus, 0, len(statusList)),
	}

	for _, status := range statusList {
		responseStatusList.JobStatusList =
			append(responseStatusList.JobStatusList, internalToExternalStatusV1(status))
	}

	return responseStatusList, nil
}

func (s *jobmanagerServer) StreamOutput(
	request *jobmanagerv1.StreamOutputRequest,
	response jobmanagerv1.JobManager_StreamOutputServer,
) error {

	userID, err := GetUserIDFromContext(response.Context())
	if err != nil {
		return err
	}

	var byteStream *io.ByteStream

	switch streamType := request.GetOutputStream(); streamType {
	case jobmanagerv1.OutputStream_OutputStream_STDOUT:
		byteStream, err = s.jm.StdoutStream(userID, request.JobID.Id)

	case jobmanagerv1.OutputStream_OutputStream_STDERR:
		byteStream, err = s.jm.StderrStream(userID, request.JobID.Id)

	default:
		return jobmanager.ErrInvalidArgument
	}

	if err != nil {
		return err
	}
	defer byteStream.Close()

	for {
		select {
		case <-response.Context().Done():
			return context.Canceled

		case data, ok := <-byteStream.Stream():
			if !ok {
				return nil
			}
			response.Send(&jobmanagerv1.JobOutput{Output: data})
		}
	}
}
