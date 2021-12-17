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

package config

import (
	"fmt"
	"os"
	"path"
)

// Note: Generally I would avoid having a "config.go" as a place for a bunch of
//       unrelated constants.  However, here these constants represent the
//       hard-coded values for the exercise, so I thought this might make it
//       easier to adjust the values for experimentation if I put them all in
//       one place.
var (
	CgexecPath string
)

// init sets CgexecPath based on the position of the current executable
func init() {
	if dir, ok := os.LookupEnv("CGEXEC_PATH"); ok {
		fmt.Println(dir)
		CgexecPath = dir
		return
	}

	const binaryName = "/cgexec"

	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	CgexecPath = path.Dir(exe) + binaryName
}

const (
	CgroupDefaultCpuLimit        = 0.5
	CgroupDefaultMemoryLimit     = "2M"
	CgroupDefaultBlkioDevice     = "8:16"
	CgroupDefaultBlkioWriteLimit = CgroupDefaultBlkioDevice + " 20971520"
	CgroupDefaultBlkioReadLimit  = CgroupDefaultBlkioDevice + " 41943040"
)
