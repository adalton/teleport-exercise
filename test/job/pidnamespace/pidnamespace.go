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
	"strings"

	"github.com/adalton/teleport-exercise/pkg/jobmanager"
)

func runTest() {

	job := jobmanager.NewJob("theOwner", "my-test", nil,
		"/bin/bash",
		"-c",
		"echo $$",
	)

	if err := job.Start(); err != nil {
		panic(err)
	}

	output := <-job.StdoutStream().Stream()
	if output == nil {
		panic("Received nil response")
	}
	outputStr := strings.TrimSpace(string(output))

	if outputStr != "1" {
		panic("The pid of the process should have been 1, not " + outputStr)
	}

	fmt.Println(outputStr)
}

// Sample run:
//     Determining the job's PID in its namespace
//     1

func main() {
	fmt.Println("Determining the job's PID in its namespace")
	runTest()
}
