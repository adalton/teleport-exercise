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

package main

import (
	"fmt"

	"github.com/adalton/teleport-exercise/pkg/cgroup/v1"
	"github.com/adalton/teleport-exercise/pkg/jobmanager"
)

func runTest(controllers ...cgroup.Controller) {

	job := jobmanager.NewJob("my-test", controllers,
		"/bin/dd",
		"if=/dev/zero",
		"of=/junk",
		"bs=4096",
		"count=100000",
		"oflag=direct",
	)

	if err := job.Start(); err != nil {
		panic(err)
	}

	for output := range job.StderrStream().Stream() {
		fmt.Print(string(output))
	}
	fmt.Println()
}

// Sample run:
//
//   $ sudo go run test/job/blkiolimit/blkiolimit.go
//   Running Blkio test with no cgroup constraints
//   100000+0 records in
//   100000+0 records out
//   409600000 bytes (410 MB, 391 MiB) copied, 2.64543 s, 155 MB/s
//
//   Running Blkio test with cgroup constraints with 8:16 20971520
//   100000+0 records in
//   100000+0 records out
//   409600000 bytes (410 MB, 391 MiB) copied, 19.5119 s, 21.0 MB/s

func main() {
	fmt.Println("Running Blkio test with no cgroup constraints")
	runTest()

	// The device portion must be a device, not a partition
	deviceString := fmt.Sprintf("8:16 %d", 1024*1024*20)
	fmt.Printf("Running Blkio test with cgroup constraints with %s\n", deviceString)

	runTest(cgroup.NewBlockIoController().
		SetReadBpsDevice(deviceString).
		SetWriteBpsDevice(deviceString))
}
