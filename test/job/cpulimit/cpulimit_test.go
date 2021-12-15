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

package cpulimit_test

import (
	"bytes"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/adalton/teleport-exercise/pkg/cgroup/cgroupv1"
	"github.com/adalton/teleport-exercise/pkg/jobmanager"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_cpulimit(t *testing.T) {
	oneCpuResult := runTest(t)
	halfCpuResult := runTest(t, &cgroupv1.CpuController{Cpus: 0.5})

	assert.True(t, aboutHalf(oneCpuResult, halfCpuResult))
}

func runTest(t *testing.T, controllers ...cgroupv1.Controller) float64 {

	job := jobmanager.NewJob("theOwner", "my-test", controllers,
		"/bin/bash",
		"-c",
		"/usr/bin/stress-ng --cpu 1 --timeout 10 --times 2>&1 | "+
			"grep 'user time' | sed -e s'/.*( *//' -e 's/%.$//'",
	)

	require.Nil(t, job.Start())

	allOutput := bytes.Buffer{}

	for output := range job.StdoutStream().Stream() {
		allOutput.Write(output)
	}

	output, err := allOutput.ReadString('\n')
	assert.Nil(t, err)

	value, err := strconv.ParseFloat(strings.TrimSpace(output), 64)
	assert.Nil(t, err)

	return value
}

func aboutHalf(firstResult, secondResult float64) bool {
	const closenessThreshold float64 = 0.5

	return math.Abs((firstResult/2.0)-secondResult) <= closenessThreshold
}
