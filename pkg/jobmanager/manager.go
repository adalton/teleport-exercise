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
	"sync"

	"github.com/adalton/teleport-exercise/pkg/cgroup/v1"
	"github.com/adalton/teleport-exercise/pkg/config"
	"github.com/adalton/teleport-exercise/pkg/io"
)

const (
	Superuser = "administrator"
)

// JobConstructor is a type that models a function for creating Jobs.
// This enables us to have a "real" job constructor function as well as
// constructor functions for mock job implementations that share the same
// signature.
type JobConstructor func(
	owner string,
	jobName string,
	controllers []cgroup.Controller,
	programPath string,
	arguments ...string,
) Job

// Manager maintains the set of jobs and enforces the authorization policy
type Manager struct {
	mutex          sync.RWMutex
	jobsByUser     map[string]map[string]Job // userId->jobId->job
	allJobs        map[string]Job            // jobId->job
	controllers    []cgroup.Controller
	jobConstructor JobConstructor
}

// NewManager creates and returns a new standard Manager.
func NewManager() *Manager {
	controllers := []cgroup.Controller{
		cgroup.NewCpuController().SetCpus(config.CgroupDefaultCpuLimit),
		cgroup.NewMemoryController().SetLimit(config.CgroupDefaultMemoryLimit),
		cgroup.NewBlockIoController().
			SetReadBpsDevice(config.CgroupDefaultBlkioReadLimit).
			SetWriteBpsDevice(config.CgroupDefaultBlkioWriteLimit),
	}

	return NewManagerDetailed(NewJob, controllers)

}

// NewManagerDetailed returns a new Manger with the given values.
// The jobConstructor is a function for creating new jobs.  In production
// this will point to NewJob.  For unit tests, this might point to a
// constructor function for a mock type.
// The given controllers is the list of cgroup controllers to manage while
// running jobs.
func NewManagerDetailed(jobConstructor JobConstructor, controllers []cgroup.Controller) *Manager {
	return &Manager{
		jobsByUser:     make(map[string]map[string]Job),
		allJobs:        make(map[string]Job),
		controllers:    controllers,
		jobConstructor: jobConstructor,
	}
}

// Start starts a new job with the given JobName for the given userId.
// The programPath and arguments are the program the user wants to run and
// the arguments to that program.
func (m *Manager) Start(userId, jobName, programPath string, arguments []string) (Job, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.jobsByUser[userId]; !exists {
		m.jobsByUser[userId] = make(map[string]Job)
	}

	job := m.jobConstructor(userId, jobName, m.controllers, programPath, arguments...)

	m.jobsByUser[userId][job.Id().String()] = job
	m.allJobs[job.Id().String()] = job

	return job, job.Start()
}

// Stop stops an existing job with the given jobId for the given userId.
func (m *Manager) Stop(userId, jobId string) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if job, err := m.findJobByUser(userId, jobId); err != nil {
		return err
	} else {
		return job.Stop()
	}
}

// List returns a list of the jobs owned by the given userId.
func (m *Manager) List(userId string) []*JobStatus {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var jobStatusList []*JobStatus

	if userId == Superuser {
		for _, job := range m.allJobs {
			jobStatusList = append(jobStatusList, job.Status())
		}
	} else {
		if l2map, exists := m.jobsByUser[userId]; exists {
			jobStatusList = make([]*JobStatus, 0, len(l2map))

			for _, job := range l2map {
				jobStatusList = append(jobStatusList, job.Status())
			}
		}
	}

	return jobStatusList
}

// Status returns the status of the job with the given JobId owned by
// the given userId.
func (m *Manager) Status(userId, jobId string) (*JobStatus, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if job, err := m.findJobByUser(userId, jobId); err != nil {
		return nil, err
	} else {
		return job.Status(), nil
	}
}

// StdoutStream returns an io.ByteStream for reading the standard output generated
// by the job with the given jobId own by the given userId.
func (m *Manager) StdoutStream(userId, jobId string) (*io.ByteStream, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if job, err := m.findJobByUser(userId, jobId); err != nil {
		return nil, err
	} else {
		return job.StdoutStream(), nil
	}

}

// Stderr returns an io.ByteStream for reading the standard error generated
// by the job with the given jobId own by the given userId.
func (m *Manager) StderrStream(userId, jobId string) (*io.ByteStream, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if job, err := m.findJobByUser(userId, jobId); err != nil {
		return nil, err
	} else {
		return job.StderrStream(), nil
	}
}

// findJobByUser finds a the job with the given jobId that is owned by
// the given userId.  If no such job is found, it returns an error.
// The caller must own the read lock associated with the given Manager.
func (m *Manager) findJobByUser(userId, jobId string) (Job, error) {
	if userId == Superuser {
		if job, exists := m.allJobs[jobId]; exists {
			return job, nil
		}
	} else {
		if l2map, exists := m.jobsByUser[userId]; exists {
			if job, exists := l2map[jobId]; exists {
				return job, nil
			}
		}
	}

	return nil, fmt.Errorf("job '%s' does not exist", jobId)
}
