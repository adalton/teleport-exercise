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

	"github.com/adalton/teleport-exercise/pkg/cgroup/cgroupv1"
	"github.com/adalton/teleport-exercise/pkg/jobmanager"
)

func runTest(controllers ...cgroupv1.Controller) {

	job := jobmanager.NewJob("theOwner", "my-test", controllers,
		"/usr/bin/stress-ng",
		"--vm",
		fmt.Sprintf("%d", runtime.NumCPU()),
		"--vm-bytes",
		fmt.Sprintf("%d", 1024*1024*1024),
		"--timeout",
		"20",
		"--oomable",
		"-v",
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
//     $ sudo go run test/job/memorylimit/memorylimit.go
//     Running Memory test with no cgroup constraints
//     stress-ng: debug: [1] stress-ng 0.13.08
//     stress-ng: debug: [1] system: Linux sinclair 5.10.76-gentoo-r1 #1 SMP Wed Nov 10 21:06:12 EST 2021 x86_64
//     stress-ng: debug: [1] RAM total: 31.3G, RAM free: 14.8G, swap free: 31.6G
//     stress-ng: debug: [1] 12 processors online, 12 processors configured
//     stress-ng: info:  [1] setting to a 20 second run per stressor
//     stress-ng: info:  [1] dispatching hogs: 12 vm
//     stress-ng: debug: [1] cache allocate: shared cache buffer size: 12288K
//     stress-ng: debug: [1] starting stressors
//     stress-ng: debug: [6] stress-ng-vm: started [6] (instance 0)
//     stress-ng: debug: [7] stress-ng-vm: started [7] (instance 1)
//     stress-ng: debug: [6] stress-ng-vm using method 'all'
//     stress-ng: debug: [8] stress-ng-vm: started [8] (instance 2)
//     stress-ng: debug: [7] stress-ng-vm using method 'all'
//     stress-ng: debug: [9] stress-ng-vm: started [9] (instance 3)
//     stress-ng: debug: [8] stress-ng-vm using method 'all'
//     stress-ng: debug: [10] stress-ng-vm: started [10] (instance 4)
//     stress-ng: debug: [9] stress-ng-vm using method 'all'
//     stress-ng: debug: [10] stress-ng-vm using method 'all'
//     stress-ng: debug: [11] stress-ng-vm: started [11] (instance 5)
//     stress-ng: debug: [11] stress-ng-vm using method 'all'
//     stress-ng: debug: [14] stress-ng-vm: started [14] (instance 6)
//     stress-ng: debug: [14] stress-ng-vm using method 'all'
//     stress-ng: debug: [16] stress-ng-vm: started [16] (instance 7)
//     stress-ng: debug: [16] stress-ng-vm using method 'all'
//     stress-ng: debug: [20] stress-ng-vm: started [20] (instance 8)
//     stress-ng: debug: [20] stress-ng-vm using method 'all'
//     stress-ng: debug: [21] stress-ng-vm: started [21] (instance 9)
//     stress-ng: debug: [1] 12 stressors started
//     stress-ng: debug: [21] stress-ng-vm using method 'all'
//     stress-ng: debug: [23] stress-ng-vm: started [23] (instance 10)
//     stress-ng: debug: [25] stress-ng-vm: started [25] (instance 11)
//     stress-ng: debug: [23] stress-ng-vm using method 'all'
//     stress-ng: debug: [25] stress-ng-vm using method 'all'
//     stress-ng: debug: [8] stress-ng-vm: exited [8] (instance 2)
//     stress-ng: debug: [7] stress-ng-vm: exited [7] (instance 1)
//     stress-ng: debug: [9] stress-ng-vm: exited [9] (instance 3)
//     stress-ng: debug: [6] stress-ng-vm: exited [6] (instance 0)
//     stress-ng: debug: [11] stress-ng-vm: exited [11] (instance 5)
//     stress-ng: debug: [14] stress-ng-vm: exited [14] (instance 6)
//     stress-ng: debug: [1] process [6] terminated
//     stress-ng: debug: [10] stress-ng-vm: exited [10] (instance 4)
//     stress-ng: debug: [1] process [7] terminated
//     stress-ng: debug: [1] process [8] terminated
//     stress-ng: debug: [1] process [9] terminated
//     stress-ng: debug: [1] process [10] terminated
//     stress-ng: debug: [1] process [11] terminated
//     stress-ng: debug: [1] process [14] terminated
//     stress-ng: debug: [23] stress-ng-vm: exited [23] (instance 10)
//     stress-ng: debug: [25] stress-ng-vm: exited [25] (instance 11)
//     stress-ng: debug: [21] stress-ng-vm: exited [21] (instance 9)
//     stress-ng: debug: [16] stress-ng-vm: exited [16] (instance 7)
//     stress-ng: debug: [1] process [16] terminated
//     stress-ng: debug: [20] stress-ng-vm: exited [20] (instance 8)
//     stress-ng: debug: [1] process [20] terminated
//     stress-ng: debug: [1] process [21] terminated
//     stress-ng: debug: [1] process [23] terminated
//     stress-ng: debug: [1] process [25] terminated
//     stress-ng: info:  [1] successful run completed in 20.05s
//     stress-ng: debug: [1] metrics-check: all stressor metrics validated and sane
//     <nil>
//
//     Running Memory test with cgroup constraints at 2M
//     stress-ng: debug: [1] stress-ng 0.13.08
//     stress-ng: debug: [1] system: Linux sinclair 5.10.76-gentoo-r1 #1 SMP Wed Nov 10 21:06:12 EST 2021 x86_64
//     stress-ng: debug: [1] RAM total: 31.3G, RAM free: 14.4G, swap free: 31.2G
//     stress-ng: debug: [1] 12 processors online, 12 processors configured
//     stress-ng: info:  [1] setting to a 20 second run per stressor
//     stress-ng: info:  [1] dispatching hogs: 12 vm
//     stress-ng: debug: [1] cache allocate: shared cache buffer size: 12288K
//     stress-ng: debug: [1] starting stressors
//     stress-ng: debug: [13] stress-ng-vm: started [13] (instance 7)
//     stress-ng: debug: [1] 12 stressors started
//     stress-ng: debug: [13] stress-ng-vm using method 'all'
//     stress-ng: debug: [12] stress-ng-vm: started [12] (instance 6)
//     stress-ng: debug: [6] stress-ng-vm: started [6] (instance 0)
//     stress-ng: debug: [9] stress-ng-vm: started [9] (instance 3)
//     stress-ng: debug: [10] stress-ng-vm: started [10] (instance 4)
//     stress-ng: debug: [8] stress-ng-vm: started [8] (instance 2)
//     stress-ng: debug: [11] stress-ng-vm: started [11] (instance 5)
//     stress-ng: debug: [7] stress-ng-vm: started [7] (instance 1)
//     stress-ng: debug: [15] stress-ng-vm: started [15] (instance 9)
//     stress-ng: debug: [14] stress-ng-vm: started [14] (instance 8)
//     stress-ng: debug: [6] stress-ng-vm using method 'all'
//     stress-ng: debug: [9] stress-ng-vm using method 'all'
//     stress-ng: debug: [12] stress-ng-vm using method 'all'
//     stress-ng: debug: [17] stress-ng-vm: started [17] (instance 11)
//     stress-ng: debug: [10] stress-ng-vm using method 'all'
//     stress-ng: debug: [15] stress-ng-vm using method 'all'
//     stress-ng: debug: [8] stress-ng-vm using method 'all'
//     stress-ng: debug: [7] stress-ng-vm using method 'all'
//     stress-ng: debug: [11] stress-ng-vm using method 'all'
//     stress-ng: debug: [16] stress-ng-vm: started [16] (instance 10)
//     stress-ng: debug: [13] stress-ng-vm: child died: signal 9 'SIGKILL' (instance 7)
//     stress-ng: debug: [14] stress-ng-vm using method 'all'
//     stress-ng: debug: [16] stress-ng-vm using method 'all'
//     stress-ng: debug: [15] stress-ng-vm: child died: signal 9 'SIGKILL' (instance 9)
//     stress-ng: debug: [6] stress-ng-vm: child died: signal 9 'SIGKILL' (instance 0)
//     stress-ng: debug: [10] stress-ng-vm: child died: signal 9 'SIGKILL' (instance 4)
//     stress-ng: debug: [13] stress-ng-vm: assuming killed by OOM killer, bailing out (instance 7)
//     stress-ng: debug: [12] stress-ng-vm: child died: signal 9 'SIGKILL' (instance 6)
//     stress-ng: debug: [6] stress-ng-vm: assuming killed by OOM killer, bailing out (instance 0)
//     stress-ng: debug: [11] stress-ng-vm: child died: signal 9 'SIGKILL' (instance 5)
//     stress-ng: debug: [15] stress-ng-vm: assuming killed by OOM killer, bailing out (instance 9)
//     stress-ng: debug: [10] stress-ng-vm: assuming killed by OOM killer, bailing out (instance 4)
//     stress-ng: debug: [12] stress-ng-vm: assuming killed by OOM killer, bailing out (instance 6)
//     stress-ng: debug: [11] stress-ng-vm: assuming killed by OOM killer, bailing out (instance 5)
//     stress-ng: debug: [13] stress-ng-vm: exited [13] (instance 7)
//     stress-ng: debug: [6] stress-ng-vm: exited [6] (instance 0)
//     stress-ng: debug: [15] stress-ng-vm: exited [15] (instance 9)
//     stress-ng: debug: [11] stress-ng-vm: exited [11] (instance 5)
//     stress-ng: debug: [12] stress-ng-vm: exited [12] (instance 6)
//     stress-ng: debug: [10] stress-ng-vm: exited [10] (instance 4)
//     stress-ng: debug: [1] process [6] terminated
//     stress-ng: debug: [1] process [7] (stress-ng-vm) terminated on signal: 9 (Killed)
//     stress-ng: debug: [1] process [7] (stress-ng-vm) was possibly killed by the OOM killer
//     stress-ng: debug: [1] process [7] terminated
//     stress-ng: debug: [1] process [8] (stress-ng-vm) terminated on signal: 9 (Killed)
//     stress-ng: debug: [1] process [8] (stress-ng-vm) was possibly killed by the OOM killer
//     stress-ng: debug: [1] process [8] terminated
//     stress-ng: debug: [1] process [9] (stress-ng-vm) terminated on signal: 9 (Killed)
//     stress-ng: debug: [1] process [9] (stress-ng-vm) was possibly killed by the OOM killer
//     stress-ng: debug: [1] process [9] terminated
//     stress-ng: debug: [1] process [10] terminated
//     stress-ng: debug: [1] process [11] terminated
//     stress-ng: debug: [1] process [12] terminated
//     stress-ng: debug: [1] process [13] terminated
//     stress-ng: debug: [16] stress-ng-vm: exited [16] (instance 10)
//     stress-ng: debug: [14] stress-ng-vm: child died: signal 9 'SIGKILL' (instance 8)
//     stress-ng: debug: [14] stress-ng-vm: assuming killed by OOM killer, bailing out (instance 8)
//     stress-ng: debug: [14] stress-ng-vm: exited [14] (instance 8)
//     stress-ng: debug: [1] process [14] terminated
//     stress-ng: debug: [1] process [15] terminated
//     stress-ng: debug: [1] process [16] terminated
//     stress-ng: debug: [1] process [17] (stress-ng-vm) terminated on signal: 9 (Killed)
//     stress-ng: debug: [1] process [17] (stress-ng-vm) was possibly killed by the OOM killer
//     stress-ng: debug: [1] process [17] terminated
//     stress-ng: info:  [1] successful run completed in 20.20s
//     stress-ng: debug: [1] metrics-check: all stressor metrics validated and sane
//     <nil>

func main() {
	fmt.Println("Running Memory test with no cgroup constraints")
	runTest()

	fmt.Println("Running Memory test with cgroup constraints at 2M")
	runTest(&cgroupv1.MemoryController{Limit: "2M"})
}
