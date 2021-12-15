//go:build integration
// +build integration

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

package pidnamespace_test

import (
	"strings"
	"testing"

	"github.com/adalton/teleport-exercise/pkg/jobmanager"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_jobPidIsOne(t *testing.T) {
	job := jobmanager.NewJob("theOwner", "my-test", nil,
		"/bin/bash",
		"-c",
		"echo $$",
	)
	defer job.Stop()

	err := job.Start()
	require.Nil(t, err)

	output := <-job.StdoutStream().Stream()
	require.NotNil(t, output)

	outputStr := strings.TrimSpace(string(output))
	assert.Equal(t, "1", outputStr)
}
