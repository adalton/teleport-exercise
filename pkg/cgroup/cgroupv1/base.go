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

import (
	"fmt"

	"github.com/adalton/teleport-exercise/pkg/adaptation/os"
)

const DefaultFileMode os.FileMode = 0644

// base is the base type for all cgroup controllers
type base struct {
	osAdapter *os.Adapter
	// name is the controller name
	name string
}

func newBase(name string, osAdapter *os.Adapter) base {
	return base{
		osAdapter: osAdapter,
		name:      name,
	}
}

// Name returns the name of this cgroup controller
func (b *base) Name() string {
	return b.name
}

// write update the cgroup file constructed with the given pathFmt and args with the given value.
func (b *base) write(value []byte, pathFmt string, args ...interface{}) error {
	return b.osAdapter.WriteFile(fmt.Sprintf(pathFmt, args...), value, DefaultFileMode)
}
