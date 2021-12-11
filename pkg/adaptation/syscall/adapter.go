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

package syscall

import gosyscall "syscall"

// Adapter serves as a shim between between callers of standard syscall.* APIs
// and the functions themselves.  The default behavior is simply to dispatch
// the call to the corresponding syscall.* function.
//
// If the field associated with that function is non-nil, then the adapter will
// dispatch to that function instead.
type Adapter struct {
	ExecFn func(argv0 string, argv []string, envv []string) (err error)
}

func (a *Adapter) Exec(argv0 string, argv []string, envv []string) (err error) {
	fn := gosyscall.Exec

	if a != nil && a.ExecFn != nil {
		fn = a.ExecFn
	}

	return fn(argv0, argv, envv)
}
