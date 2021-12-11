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
