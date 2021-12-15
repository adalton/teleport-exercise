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

	v1 "github.com/adalton/teleport-exercise/service/v1"
)

// jobmanagerServer implements the gRPC handler for the jobmanager service.
type jobmanagerServer struct {
	v1.UnimplementedJobManagerServer
}

func NewJobmanagerServer() *jobmanagerServer {
	return &jobmanagerServer{}
}

func (s *jobmanagerServer) Start(ctx context.Context, jcr *v1.JobCreationRequest) (*v1.Job, error) {
	return s.UnimplementedJobManagerServer.Start(ctx, jcr)
}

func (s *jobmanagerServer) Stop(ctx context.Context, requestJobID *v1.JobID) (*v1.NilMessage, error) {
	return s.UnimplementedJobManagerServer.Stop(ctx, requestJobID)
}

func (s *jobmanagerServer) Query(ctx context.Context, requestJobID *v1.JobID) (*v1.JobStatus, error) {
	return s.UnimplementedJobManagerServer.Query(ctx, requestJobID)
}

func (s *jobmanagerServer) List(ctx context.Context, nm *v1.NilMessage) (*v1.JobStatusList, error) {
	return s.UnimplementedJobManagerServer.List(ctx, nm)
}

func (s *jobmanagerServer) StreamOutput(
	request *v1.StreamOutputRequest,
	response v1.JobManager_StreamOutputServer,
) error {
	return s.UnimplementedJobManagerServer.StreamOutput(request, response)
}
