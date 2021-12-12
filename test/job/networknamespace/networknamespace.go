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

	"github.com/adalton/teleport-exercise/pkg/jobmanager"
)

func runTest() {

	job := jobmanager.NewJob("theOwner", "my-test", nil,
		"/bin/ip",
		"link",
	)

	if err := job.Start(); err != nil {
		panic(err)
	}

	for output := range job.StdoutStream().Stream() {
		fmt.Print(string(output))
	}
	fmt.Printf("\n")
}

// Sample run:
//     $ sudo go run networknamespace.go
//     Running test to list all network interfaces avaialble to a job
//     1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN mode DEFAULT group default qlen 1000
//         link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
//     2: sit0@NONE: <NOARP> mtu 1480 qdisc noop state DOWN mode DEFAULT group default qlen 1000
//         link/sit 0.0.0.0 brd 0.0.0.0

func main() {
	fmt.Println("Running test to list all network interfaces avaialble to a job")
	runTest()
}
