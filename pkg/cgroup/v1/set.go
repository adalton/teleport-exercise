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

package cgroup

import (
	"fmt"
	"log"
	"strings"

	"github.com/adalton/teleport-exercise/pkg/adaptation/os"

	"github.com/google/uuid"
)

const (
	DefaultBasePath                   = "/sys/fs/cgroup"
	defaultDirectoryPerms os.FileMode = 0755
)

// Set maintains a collection of 0 or more cgroup controllers that should be
// created/removed at the same time.
type Set struct {
	osAdapter   *os.Adapter
	basePath    string
	jobId       uuid.UUID
	controllers []Controller
}

// NewSet creates a new cgroup (v1) set for the given jobID.  This assumes that
// the cgroup filesystem is mounted at /sys/fs/cgroup.
func NewSet(jobId uuid.UUID, controllers ...Controller) *Set {
	return NewSetDetailed(nil, DefaultBasePath, jobId, controllers...)
}

// NewSetDetailed creates a new cgroup (v1) set for the given jobID rooted
// at the given basePath.
func NewSetDetailed(
	osAdapter *os.Adapter,
	basePath string,
	jobId uuid.UUID,
	controllers ...Controller,
) *Set {

	return &Set{
		osAdapter:   osAdapter,
		basePath:    basePath,
		jobId:       jobId,
		controllers: controllers,
	}
}

// Create creates the cgroup v1 directories for all registered controllers.
func (s *Set) Create() (retErr error) {
	if s == nil {
		// If the set is nil, then Create is vacuously successful
		return nil
	}

	failPoint := -1

	for i := range s.controllers {
		path := fmt.Sprintf("%s/%s/jobs/%s", s.basePath, s.controllers[i].Name(), s.jobId.String())

		// Create the cgroup
		if err := s.osAdapter.MkdirAll(path, defaultDirectoryPerms); err != nil {
			retErr = err
			failPoint = i
			break
		}

		// Apply the supplied configuration to the newly-created cgroup
		if err := s.controllers[i].Apply(path); err != nil {
			retErr = err
			failPoint = i
			break
		}
	}

	for i := failPoint; i >= 0; i-- {
		path := fmt.Sprintf("%s/%s/jobs/%s", s.basePath, s.controllers[i].Name(), s.jobId.String())

		if err := s.osAdapter.Remove(path); err != nil {
			_ = err
			log.Printf("Failed to backout cgroup %s: %v", s.controllers[i].Name(), err)
			// Intentionally not returning an error here
		}
	}

	return
}

// Destroy removes the cgroup v1 directories for all registered controllers.
func (s *Set) Destroy() error {
	if s == nil {
		// If the set is nil, then Delete is vacuously successful
		return nil
	}

	var failedCgroups []string

	for i := len(s.controllers) - 1; i >= 0; i-- {
		path := fmt.Sprintf("%s/%s/jobs/%s", s.basePath, s.controllers[i].Name(), s.jobId.String())

		if err := s.osAdapter.Remove(path); err != nil {
			failedCgroups = append(failedCgroups, path)
		}
	}

	if len(failedCgroups) > 0 {
		return fmt.Errorf("failed to destroy cgroups: %s", strings.Join(failedCgroups, ", "))
	}

	return nil
}

// TaskFiles returns a list of the 'tasks' files for each cgroup controller
// in this set.
func (s *Set) TaskFiles() []string {
	if s == nil {
		return nil
	}

	taskFiles := make([]string, 0, len(s.controllers))

	for i := range s.controllers {
		taskFiles = append(taskFiles, fmt.Sprintf(
			"%s/%s/jobs/%s/tasks",
			s.basePath,
			s.controllers[i].Name(),
			s.jobId.String()))
	}

	return taskFiles
}
