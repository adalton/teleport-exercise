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
	"os"

	"github.com/adalton/teleport-exercise/pkg/command"
)

// main is the entrypoint of the cgexec application.  The application accepts
// arguments in the form:
//
//     cgexec [<cgtskfile> ...] -- <command> [<arg> ...]
//
// Everything before the first "--" is treated as a cgroup task file; this
// command will add itself to those cgroups.
//
// Everything after the first "--" is treated as the command and arguments to
// the program to exec.
//
// If no "--" is found, then all arguments are treated as the command and
// arguments to the program to exec.
func main() {
	if err := command.Cgexec(os.Args); err != nil {
		fmt.Printf("cgexec failed: %v", err)
	}

	// Cgexec shouldn't return in a non-error case
	os.Exit(1)
}
