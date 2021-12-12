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
	"github.com/adalton/teleport-exercise/pkg/io"
)

const (
	Superuser = "administrator"
)

type JobConstructor func(
	jobName string,
	controllers []cgroup.Controller,
	programPath string,
	arguments ...string,
) Job

type Manager struct {
	mutex          sync.RWMutex
	jobsByUser     map[string]map[string]Job // userId->jobId->job
	allJobs        map[string]Job            // jobId->job
	controllers    []cgroup.Controller
	jobConstructor JobConstructor
}

func NewManager() *Manager {
	// TODO:
	readLimit := fmt.Sprintf("8:16 %d", 1024*1024*20)
	writeLimit := fmt.Sprintf("8:16 %d", 1024*1024*40)

	controllers := []cgroup.Controller{
		cgroup.NewCpuController().SetCpus(0.5),
		cgroup.NewMemoryController().SetLimit("2M"),
		cgroup.NewBlockIoController().
			SetReadBpsDevice(readLimit).
			SetWriteBpsDevice(writeLimit),
	}

	return NewManagerDetailed(NewJob, controllers)

}

func NewManagerDetailed(jobConstructor JobConstructor, controllers []cgroup.Controller) *Manager {
	return &Manager{
		jobsByUser:     make(map[string]map[string]Job),
		allJobs:        make(map[string]Job),
		controllers:    controllers,
		jobConstructor: jobConstructor,
	}
}

func (m *Manager) Start(userId, jobName, programPath string, arguments []string) (Job, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.jobsByUser[userId]; !exists {
		m.jobsByUser[userId] = make(map[string]Job)
	}

	job := m.jobConstructor(jobName, m.controllers, programPath, arguments...)

	m.jobsByUser[userId][job.Id().String()] = job
	m.allJobs[job.Id().String()] = job

	return job, job.Start()
}

func (m *Manager) Stop(userId, jobId string) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if job, err := m.findJobByUser(userId, jobId); err != nil {
		return err
	} else {
		return job.Stop()
	}
}

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

func (m *Manager) Status(userId, jobId string) (*JobStatus, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if job, err := m.findJobByUser(userId, jobId); err != nil {
		return nil, err
	} else {
		return job.Status(), nil
	}
}

func (m *Manager) StdoutStream(userId, jobId string) (*io.ByteStream, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if job, err := m.findJobByUser(userId, jobId); err != nil {
		return nil, err
	} else {
		return job.StdoutStream(), nil
	}

}

func (m *Manager) StderrStream(userId, jobId string) (*io.ByteStream, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if job, err := m.findJobByUser(userId, jobId); err != nil {
		return nil, err
	} else {
		return job.StderrStream(), nil
	}
}

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
