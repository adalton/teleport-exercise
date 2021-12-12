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
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/adalton/teleport-exercise/pkg/cgroup/v1"
	"github.com/adalton/teleport-exercise/pkg/config"
	"github.com/adalton/teleport-exercise/pkg/io"
	"github.com/google/uuid"
)

// JobStatus models the current status of a job.
type JobStatus struct {
	Owner     string
	Name      string
	Id        string
	Running   bool
	Pid       int
	ExitCode  int
	SignalNum syscall.Signal
	RunError  error
}

// Job defines an interface for objects that behave like jobs.  This enables
// us to define both a production job type as well a a job type for unit
// testing.
type Job interface {
	Start() error
	Stop() error
	Status() *JobStatus
	StdoutStream() *io.ByteStream
	StderrStream() *io.ByteStream
	Name() string
	Id() uuid.UUID
}

// concreteJob implements the Job interface and provides the production implementation
// of concreteJob behavior.
type concreteJob struct {
	mutex         sync.Mutex
	owner         string
	id            uuid.UUID
	name          string
	cgControllers []cgroup.Controller
	programName   string
	programArgs   []string
	cmd           *exec.Cmd
	stdoutBuffer  io.OutputBuffer
	stderrBuffer  io.OutputBuffer
	running       bool
	runErrors     []error
}

// NewJob creates and returns a new concreteJob based on the given values.
func NewJob(
	owner string,
	name string,
	cgControllers []cgroup.Controller,
	programName string,
	programArgs ...string,
) Job {

	return NewJobDetailed(
		owner,
		name,
		cgControllers,
		io.NewMemoryBuffer(),
		io.NewMemoryBuffer(),
		programName,
		programArgs...,
	)
}

func NewJobDetailed(
	owner string,
	name string,
	cgControllers []cgroup.Controller,
	stdoutBuffer io.OutputBuffer,
	stderrBuffer io.OutputBuffer,
	programName string,
	programArgs ...string,
) Job {

	return &concreteJob{
		owner:         owner,
		id:            uuid.New(),
		name:          name,
		cgControllers: cgControllers,
		programName:   programName,
		programArgs:   programArgs,
		stdoutBuffer:  stdoutBuffer,
		stderrBuffer:  stderrBuffer,
	}
}

// Start starts the job if it hasn't already been started.
func (j *concreteJob) Start() error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.cmd != nil {
		return fmt.Errorf("job %s (%v) has already been started", j.name, j.id)
	}

	cgroupSet := cgroup.NewSet(j.id, j.cgControllers...)
	if err := cgroupSet.Create(); err != nil {
		return err
	}

	args := cgroupSet.TaskFiles()
	args = append(args, "--")
	args = append(args, j.programName)
	args = append(args, j.programArgs...)

	j.cmd = exec.Command(config.CgexecPath, args...)
	j.cmd.Stdout = j.stdoutBuffer
	j.cmd.Stderr = j.stderrBuffer
	j.cmd.Env = make([]string, 0) // Do not pass along our environment

	//j.cmd.Dir = "/" // If we were to chroot

	j.cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: "", // This would be non-empty to actually do a chroot
		Cloneflags: syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET,
	}

	j.running = true

	go func() {
		defer j.lockedOperation(func() {
			if err := cgroupSet.Destroy(); err != nil {
				j.runErrors = append(j.runErrors, err)
			}
			j.running = false
		})

		// Run blocks until the newly-created process terminates.  It calls
		// Wait internally
		err := j.cmd.Run()

		// Once Wait returns, all output has been written to Stdout and Stderr
		j.lockedOperation(func() {
			if err != nil {
				j.runErrors = append(j.runErrors, err)
			}

			if err := j.stdoutBuffer.Close(); err != nil {
				j.runErrors = append(j.runErrors, err)
			}

			if err := j.stderrBuffer.Close(); err != nil {
				j.runErrors = append(j.runErrors, err)
			}
		})
	}()

	return nil
}

// Stop kills the job.
func (j *concreteJob) Stop() error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if !j.running {
		// If the job isn't running, it is stopped already
		return nil
	}

	if err := j.cmd.Process.Kill(); err != nil && err != os.ErrProcessDone {
		return err
	}

	return nil
}

// StdoutStream returns a ByteStream associated with the standard output of the job.
func (j *concreteJob) StdoutStream() *io.ByteStream {
	// Unlocked read of j.stdoutBuffer should be ok since it's not modified once created
	return io.NewByteStream(j.stdoutBuffer)
}

// StderrStream returns a ByteStream associated with the standard error of the job.
func (j *concreteJob) StderrStream() *io.ByteStream {
	// Unlocked read of j.stderrBuffer should be ok since it's not modified once created
	return io.NewByteStream(j.stderrBuffer)
}

// Status returns the current status of this job.  If the job is running,
// the information will include the job's PID.  If the job has terminated,
// the information will include the exit code and termination signal (if any).
func (j *concreteJob) Status() *JobStatus {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	status := &JobStatus{
		Owner:     j.owner,
		Name:      j.name,
		Id:        j.id.String(),
		Running:   j.running,
		Pid:       -1,
		SignalNum: syscall.Signal(-1),
	}

	if j.runErrors != nil {
		status.RunError = fmt.Errorf("run error: %v", j.runErrors)
	}

	if state := j.cmd.ProcessState; state == nil {
		if j.cmd.Process != nil {
			status.Pid = j.cmd.Process.Pid
		}
	} else {
		if sys := state.Sys(); sys != nil {
			status.SignalNum = sys.(syscall.WaitStatus).Signal()
		}
	}

	return status
}

// Id returns the server-assigned ID of this job.
func (j *concreteJob) Id() uuid.UUID {
	return j.id
}

// Name returns the user-assigned name of this job.
func (j *concreteJob) Name() string {
	return j.name
}

// lockedOperation is a simple runs the functor with the concreteJob lock held.
// The caller must not hold the lock.
func (j *concreteJob) lockedOperation(fn func()) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	fn()
}
