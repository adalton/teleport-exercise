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

package blkiolimit_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/adalton/teleport-exercise/pkg/cgroup/cgroupv1"
	"github.com/adalton/teleport-exercise/pkg/jobmanager"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Todo: writing to the root filesystem isn't ideal.  In a real scenario
	//       this would be configurable.
	tmpFileDirectory = "/"

	deviceString = "8:16"
)

func Test_blkiolimit(t *testing.T) {
	noLimit := runTest(t)

	// The device portion must be a device, not a partition
	deviceString := fmt.Sprintf("%s %d", deviceString, 1024*1024*20)
	withLimit := runTest(t, &cgroupv1.BlockIOController{
		ReadBpsDevice:  deviceString,
		WriteBpsDevice: deviceString,
	})

	// Give it a little wiggle room.  This might need some additional experimentation
	// to dial in on a suitable value.
	const limitThreshold = 2.0

	assert.Less(t, withLimit, noLimit)
	assert.Less(t, withLimit, 20.0+limitThreshold)
}

func runTest(t *testing.T, controllers ...cgroupv1.Controller) float64 {

	file, err := ioutil.TempFile(tmpFileDirectory, "blkiolimit-test")
	require.Nil(t, err)
	defer os.Remove(file.Name())

	job := jobmanager.NewJob("theOwner", "my-test", controllers,
		"/bin/bash",
		"-c",
		"/bin/dd if=/dev/zero of="+file.Name()+" bs=4096 count=100000 oflag=direct 2>&1 |"+
			"grep copied | sed -e 's-^.*, --' -e 's/ .*$//'",
	)

	require.Nil(t, job.Start())

	allOutput := bytes.Buffer{}
	for output := range job.StdoutStream().Stream() {
		allOutput.Write(output)
	}

	output, err := allOutput.ReadString('\n')
	assert.Nil(t, err)

	value, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSpace(output)), 64)
	assert.Nil(t, err)

	return value
}
