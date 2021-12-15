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

package syscalltest

import "fmt"

// ExecMock is a mock implementation of the Exec system call wrapper.
// This implementation returns the received parameters and returns the
// configured error. If no error is specified, this implementation will
// create one to return.
type ExecMock struct {
	Argv0 string
	Argv  []string
	Envv  []string
	Error error
}

func (n *ExecMock) Exec(argv0 string, argv []string, envv []string) (err error) {
	n.Argv0 = argv0
	n.Argv = argv
	n.Envv = envv

	err = n.Error
	if err == nil {
		// Exec must never return a non-error value
		err = fmt.Errorf("nilexec: exec failed")
	}

	return err
}
