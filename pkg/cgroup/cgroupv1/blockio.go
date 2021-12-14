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

const (
	BlkioThrottleReadBpsDevice  = "blkio.throttle.read_bps_device"
	BlkioThrottleWriteBpsDevice = "blkio.throttle.write_bps_device"
)

// BlockIOController implements the BlockIOController cgroup controller
type BlockIOController struct {
	OsAdapter      *os.Adapter
	ReadBpsDevice  string
	WriteBpsDevice string
}

func (BlockIOController) Name() string {
	return "blkio"
}

// Apply applies this cgroup controller configuration to the blkio cgroup
// at the given path.
func (b *BlockIOController) Apply(path string) error {
	if b.ReadBpsDevice != "" {
		filename := fmt.Sprintf("%s/%s", path, BlkioThrottleReadBpsDevice)

		if err := b.OsAdapter.WriteFile(filename, []byte(b.ReadBpsDevice), os.FileMode(0644)); err != nil {
			return err
		}
	}

	if b.WriteBpsDevice != "" {
		filename := fmt.Sprintf("%s/%s", path, BlkioThrottleWriteBpsDevice)

		if err := b.OsAdapter.WriteFile(filename, []byte(b.WriteBpsDevice), os.FileMode(0644)); err != nil {
			return err
		}
	}

	return nil
}
