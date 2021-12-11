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

package cgroup

import "github.com/adalton/teleport-exercise/pkg/adaptation/os"

const (
	BlkioThrottleReadBpsDevice  = "blkio.throttle.read_bps_device"
	BlkioThrottleWriteBpsDevice = "blkio.throttle.write_bps_device"
)

type blockIo struct {
	base
	// TODO: It might be helpful to enable these to be lists so that a single
	// blocIo object can handle multiple devices
	readBpsDevice  *string
	writeBpsDevice *string
}

func NewBlockIoController() *blockIo {
	return NewBlockIoControllerDetailed(nil)
}

func NewBlockIoControllerDetailed(osAdapter *os.Adapter) *blockIo {
	return &blockIo{
		base: newBase("blkio", osAdapter),
	}
}

func (b *blockIo) Apply(path string) error {
	if b.readBpsDevice != nil {
		if err := b.write([]byte(*b.readBpsDevice), "%s/%s", path, BlkioThrottleReadBpsDevice); err != nil {
			return err
		}
	}

	if b.writeBpsDevice != nil {
		if err := b.write([]byte(*b.writeBpsDevice), "%s/%s", path, BlkioThrottleWriteBpsDevice); err != nil {
			return err
		}
	}

	return nil
}

func (b *blockIo) SetReadBpsDevice(value string) *blockIo {
	b.readBpsDevice = &value

	return b
}

func (b *blockIo) SetWriteBpsDevice(value string) *blockIo {
	b.writeBpsDevice = &value

	return b
}
