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
		"/usr/bin/id",
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
//     Determining the job's PID in its namespace
//     uid=0(root) gid=0(root) groups=0(root)

func main() {
	fmt.Println("Determining the job's PID in its namespace")
	runTest()
}
