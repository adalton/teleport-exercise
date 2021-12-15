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

package cgroupv1

// Controller defines the interface to a cgroup controller -- objects that
// model concrete cgroup controlers and their configuration options.
type Controller interface {
	// Name returns the name of the cgroup
	Name() string

	// Apply applies this controller's configuration to the cgroup at the
	// given path.
	Apply(path string) error
}
