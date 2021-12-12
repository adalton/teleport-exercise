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
	"runtime"

	"github.com/adalton/teleport-exercise/pkg/cgroup/v1"
	"github.com/adalton/teleport-exercise/pkg/jobmanager"
)

func runTest(controllers ...cgroup.Controller) {

	job := jobmanager.NewJob("my-test", controllers, "/usr/bin/stress-ng",
		"--cpu",
		fmt.Sprintf("%d", runtime.NumCPU()),
		"--timeout",
		"20",
		"--times",
	)

	if err := job.Start(); err != nil {
		panic(err)
	}

	for output := range job.StderrStream().Stream() {
		fmt.Print(string(output))
	}
	fmt.Printf("\n")
}

// Sample run:
//     $ sudo go run test/job/cpulimit/cpulimit.go
//     Password:
//     Running CPU test with no cgroup constraints
//     stress-ng: info:  [1] setting to a 20 second run per stressor
//     stress-ng: info:  [1] dispatching hogs: 12 cpu
//     stress-ng: info:  [1] successful run completed in 20.01s
//     stress-ng: info:  [1] for a 20.01s run time:
//     stress-ng: info:  [1]     240.09s available CPU time
//     stress-ng: info:  [1]     239.37s user time   ( 99.70%)
//     stress-ng: info:  [1]       0.05s system time (  0.02%)
//     stress-ng: info:  [1]     239.42s total time  ( 99.72%)
//     stress-ng: info:  [1] load average: 5.83 2.88 2.22
//
//     Running CPU test with cgroup constraints at 0.5 CPU
//     stress-ng: info:  [1] setting to a 20 second run per stressor
//     stress-ng: info:  [1] dispatching hogs: 12 cpu
//     stress-ng: info:  [1] successful run completed in 20.13s
//     stress-ng: info:  [1] for a 20.13s run time:
//     stress-ng: info:  [1]     241.60s available CPU time
//     stress-ng: info:  [1]      10.08s user time   (  4.17%)
//     stress-ng: info:  [1]       0.02s system time (  0.01%)
//     stress-ng: info:  [1]      10.10s total time  (  4.18%)
//     stress-ng: info:  [1] load average: 5.47 2.99 2.27

func main() {
	fmt.Println("Running CPU test with no cgroup constraints")
	runTest()

	fmt.Println("Running CPU test with cgroup constraints at 0.5 CPU")
	runTest(cgroup.NewCpuController().SetCpus(0.5))
}
