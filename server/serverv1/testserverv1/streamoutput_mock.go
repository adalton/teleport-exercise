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

package testserverv1

import (
	"context"

	v1 "github.com/adalton/teleport-exercise/service/v1"
	"google.golang.org/grpc/metadata"
)

// MockJobManagerStreamServer mocks the APIs used by a JobManager server.
type MockJobmanagerStreamServer struct {
	LastJobOutput *v1.JobOutput
	SendCount     int
	SendError     error
	NextContext   context.Context
}

func (m *MockJobmanagerStreamServer) Send(output *v1.JobOutput) error {
	m.SendCount++
	m.LastJobOutput = output
	return m.SendError
}

func (m *MockJobmanagerStreamServer) Context() context.Context {
	return m.NextContext
}

// SetHeader is not yet implemented; it will panic.
func (m *MockJobmanagerStreamServer) SetHeader(metadata.MD) error {
	panic("unimplemented")
}

// SendHeader is not yet implemented; it will panic.
func (m *MockJobmanagerStreamServer) SendHeader(metadata.MD) error {
	panic("unimplemented")
}

// SetTrailer is not yet implemented; it will panic.
func (m *MockJobmanagerStreamServer) SetTrailer(metadata.MD) {
	panic("unimplemented")
}

// SendMsg is not yet implemented; it will panic.
func (m *MockJobmanagerStreamServer) SendMsg(interface{}) error {
	panic("unimplemented")
}

// RecvMsg is not yet implemented; it will panic.
func (m *MockJobmanagerStreamServer) RecvMsg(interface{}) error {
	panic("unimplemented")
}
