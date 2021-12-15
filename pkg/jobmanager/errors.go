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
	"errors"
	"fmt"
)

type JobExistsError struct {
	jobName string
}

func NewJobExistsError(jobName string) error {
	return &JobExistsError{jobName: jobName}
}

func (j *JobExistsError) Error() string {
	return fmt.Sprintf("job with name '%s' already exists", j.jobName)
}

func (j *JobExistsError) Is(err error) bool {
	_, ok := err.(*JobExistsError)

	return err != nil && ok
}

type JobNotFoundError struct {
	jobID string
}

func NewJobNotFoundError(jobID string) error {
	return &JobNotFoundError{jobID: jobID}
}

func (j *JobNotFoundError) Error() string {
	return fmt.Sprintf("job with ID '%s' not found", j.jobID)
}

func (j *JobNotFoundError) Is(err error) bool {
	_, ok := err.(*JobNotFoundError)

	errors.Is(nil, nil)

	return err != nil && ok
}

type InvalidJobID struct {
	jobID string
}

func NewInvalidJobID(jobID string) error {
	return &InvalidJobID{jobID: jobID}
}

func (j *InvalidJobID) Error() string {
	return fmt.Sprintf("'%s' is not a valid job ID", j.jobID)
}

func (j *InvalidJobID) Is(err error) bool {
	_, ok := err.(*InvalidJobID)

	errors.Is(nil, nil)

	return err != nil && ok
}
