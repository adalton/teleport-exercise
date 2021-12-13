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
	"sync"

	"github.com/adalton/teleport-exercise/pkg/jobmanager"
)

func runTest() {

	job := jobmanager.NewJob("theOwner", "my-test", nil,
		"/bin/bash",
		"-c",
		"for ((i = 0; i < 100; ++i)); do for((j = 0; j < 1000; ++j)); do echo $RANDOM; done; sleep 0.25; done",
	)

	if err := job.Start(); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(threadNum int) {
			count := 0
			for output := range job.StdoutStream().Stream() {
				count += len(output)
				if threadNum == 0 {
					fmt.Print(string(output))
				}
			}
			fmt.Printf("%d: %d\n", threadNum, count)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

// The job generates 10000 random numbers and prints them to standard output
// The program starts 100 goroutines to consume that output.  Each goroutine
// counts and prints the number of bytes that it receives.

func main() {
	runTest()
}
