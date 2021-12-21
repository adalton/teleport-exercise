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

package jobmanagertest

import (
	"fmt"
	"syscall"

	"github.com/adalton/teleport-exercise/pkg/cgroup/cgroupv1"
	"github.com/adalton/teleport-exercise/pkg/io"
	"github.com/adalton/teleport-exercise/pkg/jobmanager"

	"github.com/google/uuid"
)

const (
	DefaultStandardOutput         = "this is standard output"
	DefaultStandardError          = "this is standard error"
	DefaultPID                    = 1234
	DefaultSignalWhileRunning     = syscall.Signal(-1)
	DefaultSignalAfterStop        = syscall.SIGKILL
	DefaultExitStatusWhileRunning = -1
	DefaultExitStatusAfterStop    = 128 + int(DefaultSignalAfterStop)
)

// mockJob is a simple implementation of the Job interface for use by unit tests
type mockJob struct {
	owner   string
	name    string
	id      uuid.UUID
	running bool
	stdout  io.OutputBuffer
	stderr  io.OutputBuffer
}

// NewMockJob creates and returns a new mockJob.
func NewMockJob(
	owner string,
	jobName string,
	controllers []cgroupv1.Controller,
	programPath string,
	arguments ...string,
) jobmanager.Job {
	return &mockJob{
		owner: owner,
		name:  jobName,
		// Normally I'd using random values in a unit test, but here I wanted
		// this constructor to match the signature of the one for concreteJob.
		id:     uuid.New(),
		stdout: io.NewMemoryBuffer(),
		stderr: io.NewMemoryBuffer(),
	}
}

func (m *mockJob) Start() error {

	if m.running {
		return fmt.Errorf("job %s (%v) has already been started", m.name, m.id)
	}

	m.running = true
	_, _ = m.stdout.Write([]byte(DefaultStandardOutput))
	_, _ = m.stderr.Write([]byte(DefaultStandardError))

	return nil
}

func (m *mockJob) Stop() error {
	m.running = false
	m.stdout.Close()
	m.stderr.Close()
	return nil
}

func (m *mockJob) StdoutStream() *io.ByteStream {
	return io.NewByteStream(m.stdout)
}

func (m *mockJob) StderrStream() *io.ByteStream {
	return io.NewByteStream(m.stderr)
}

func (m *mockJob) Status() *jobmanager.JobStatus {
	exitCode := DefaultExitStatusWhileRunning
	signalNumber := DefaultSignalWhileRunning

	if !m.running {
		exitCode = DefaultExitStatusAfterStop
		signalNumber = DefaultSignalAfterStop
	}

	return &jobmanager.JobStatus{
		Owner:     m.owner,
		Name:      m.name,
		ID:        m.id.String(),
		Running:   m.running,
		Pid:       DefaultPID,
		SignalNum: signalNumber,
		ExitCode:  exitCode,
		RunError:  nil,
	}
}

func (m *mockJob) ID() uuid.UUID {
	return m.id
}

func (m *mockJob) Name() string {
	return m.name
}
