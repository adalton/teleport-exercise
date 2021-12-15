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
	v1 "github.com/adalton/teleport-exercise/service/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserIDContext struct{}

// jobmanagerServer implements the gRPC handler for the jobmanager service.
type jobmanagerServer struct {
	v1.UnimplementedJobManagerServer
	jm *jobmanager.Manager
}

// NewJobmanagerServer creates and returns a jobmanagerServer.
func NewJobmanagerServer() *jobmanagerServer {
	return &jobmanagerServer{
		jm: jobmanager.NewManager(),
	}
}

func (s *jobmanagerServer) Start(ctx context.Context, jcr *v1.JobCreationRequest) (*v1.Job, error) {
	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	job, err := s.jm.Start(userID, jcr.GetName(), jcr.GetProgramPath(), jcr.GetArguments())
	if err != nil {
		// Todo: Examine Start implementation.  See if there are any errors
		//       that we should map to more meaningful gRCP error codes.
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	jobResponse := &v1.Job{
		Id:   &v1.JobID{Id: job.ID().String()},
		Name: job.Name(),
	}

	return jobResponse, nil
}

func (s *jobmanagerServer) Stop(ctx context.Context, requestJobID *v1.JobID) (*v1.NilMessage, error) {
	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = s.jm.Stop(userID, requestJobID.Id)
	if err != nil {
		// Todo: Examine Stop implementation.  See if there are any errors
		//       that we should map to more meaningful gRCP error codes.
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &v1.NilMessage{}, nil
}

func internalToExternalStatusV1(internalStatus *jobmanager.JobStatus) *v1.JobStatus {
	errMsg := ""

	if internalStatus.RunError != nil {
		errMsg = internalStatus.RunError.Error()
	}

	return &v1.JobStatus{
		Job: &v1.Job{
			Id: &v1.JobID{
				Id: internalStatus.ID,
			},
			Name: internalStatus.Name,
		},
		IsRunning:    internalStatus.Running,
		ExitCode:     int32(internalStatus.ExitCode),
		SignalNumber: int32(internalStatus.SignalNum),
		ErrorMessage: errMsg,
	}
}

func (s *jobmanagerServer) Query(ctx context.Context, requestJobID *v1.JobID) (*v1.JobStatus, error) {
	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	jobStatus, err := s.jm.Status(userID, requestJobID.Id)
	if err != nil {
		// Todo: Examine Status implementation.  See if there are any errors
		//       that we should map to more meaningful gRCP error codes.
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return internalToExternalStatusV1(jobStatus), nil
}

func (s *jobmanagerServer) List(ctx context.Context, _ *v1.NilMessage) (*v1.JobStatusList, error) {
	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	statusList := s.jm.List(userID)

	responseStatusList := &v1.JobStatusList{
		JobStatusList: make([]*v1.JobStatus, 0, len(statusList)),
	}

	for _, status := range statusList {
		responseStatusList.JobStatusList =
			append(responseStatusList.JobStatusList, internalToExternalStatusV1(status))
	}

	return responseStatusList, nil
}

func (s *jobmanagerServer) StreamOutput(
	request *v1.StreamOutputRequest,
	response v1.JobManager_StreamOutputServer,
) error {

	userID, err := s.userIDFromContext(response.Context())
	if err != nil {
		return err
	}

	var byteStream *io.ByteStream

	switch streamType := request.GetOutputStream(); streamType {
	case v1.OutputStream_STDOUT:
		byteStream, err = s.jm.StdoutStream(userID, request.JobID.Id)

	case v1.OutputStream_STDERR:
		byteStream, err = s.jm.StderrStream(userID, request.JobID.Id)

	default:
		return status.Errorf(codes.InvalidArgument, "unsupported stream: %v", streamType)
	}

	if err != nil {
		// Todo: Examine Stream* implementation.  See if there are any errors
		//       that we should map to more meaningful gRCP error codes.
		return status.Errorf(codes.Unknown, err.Error())
	}

	for data := range byteStream.Stream() {
		response.Send(&v1.JobOutput{Output: data})
	}

	return nil
}

func (s *jobmanagerServer) userIDFromContext(ctx context.Context) (string, error) {
	if userID, ok := ctx.Value(&UserIDContext{}).(string); ok && userID != "" {
		return userID, nil
	}

	return "", status.Errorf(codes.Unauthenticated, "jobmanager: user unauthenticated")
}
