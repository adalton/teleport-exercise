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
	"encoding/json"
	"fmt"

	"github.com/adalton/teleport-exercise/pkg/jobmanager"
)

func runTest() {

	job := jobmanager.NewJob("theOwner", "my-test", nil,
		"/bin/ip",
		"-j",
		"link",
	)

	if err := job.Start(); err != nil {
		panic(err)
	}

	var outputBuffer []byte

	for output := range job.StdoutStream().Stream() {
		outputBuffer = append(outputBuffer, output...)
	}

	type iface struct {
		Ifindex *int    `json:"ifindex,omitempty"`
		Ifname  *string `json:"ifname,omitempty"`
	}
	var ifaceList []iface

	if err := json.Unmarshal(outputBuffer, &ifaceList); err != nil {
		panic(err)
	}

	if len(ifaceList) != 2 {
		panic(fmt.Sprintf("Expected 2, found: %d", len(ifaceList)))
	}

	fmt.Println("Found expected number of network interface in new network namespace (2)")
}

// Sample run:
//     Running test to list all network interfaces avaialble to a job
//     Found expected number of network interface in new network namespace (2)

func main() {
	fmt.Println("Running test to list all network interfaces avaialble to a job")
	runTest()
}
