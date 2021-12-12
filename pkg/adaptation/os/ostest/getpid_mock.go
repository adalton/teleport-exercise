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

package ostest

// GetpidMock is a component that provides a mock implementation of the
// os.Getpid() function.  The function returns the configured Pid.
type GetpidMock struct {
	Pid int
}

func (p *GetpidMock) Getpid() int {
	if p.Pid == 0 {
		return 1
	}

	return p.Pid
}