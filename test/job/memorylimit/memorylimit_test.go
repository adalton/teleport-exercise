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

package memorylimit_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/adalton/teleport-exercise/pkg/cgroup/cgroupv1"
	"github.com/adalton/teleport-exercise/pkg/jobmanager"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_memorylimit(t *testing.T) {
	noLimitCount := runTest(t)
	require.Equal(t, 0, noLimitCount)

	limitCount := runTest(t, &cgroupv1.MemoryController{Limit: "1M"})
	require.Greater(t, limitCount, 0)
}

func runTest(t *testing.T, controllers ...cgroupv1.Controller) int {

	cmd := fmt.Sprintf("/usr/bin/stress-ng --vm 1 --vm-bytes %d --timeout 10 --oomable -v 2>&1 | grep 'OOM killer'", 1024*1024*1024)

	job := jobmanager.NewJob("theOwner", "my-test", controllers,
		"/bin/bash",
		"-c",
		cmd)

	require.Nil(t, job.Start())

	outputBuffer := bytes.Buffer{}
	for output := range job.StdoutStream().Stream() {
		outputBuffer.Write(output)
	}

	lineCount := 0

	var err error

	for {
		_, err = outputBuffer.ReadString('\n')
		if err != nil {
			assert.Equal(t, err, io.EOF)
			break
		}

		lineCount++
	}

	return lineCount
}
