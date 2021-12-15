//go:build integration
// +build integration

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

package networknamespace_test

import (
	"encoding/json"
	"testing"

	"github.com/adalton/teleport-exercise/pkg/jobmanager"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_networknamespace(t *testing.T) {
	job := jobmanager.NewJob("theOwner", "my-test", nil,
		"/bin/ip",
		"-j",
		"link",
	)

	require.Nil(t, job.Start())
	defer job.Stop()

	var outputBuffer []byte

	for output := range job.StdoutStream().Stream() {
		outputBuffer = append(outputBuffer, output...)
	}

	type iface struct {
		Ifname *string `json:"ifname,omitempty"`
	}
	var ifaceList []iface

	err := json.Unmarshal(outputBuffer, &ifaceList)
	assert.Nil(t, err)

	require.Equal(t, 2, len(ifaceList))
	require.NotNil(t, ifaceList[0].Ifname)
	require.NotNil(t, ifaceList[1].Ifname)
	assert.Equal(t, "lo", *ifaceList[0].Ifname)
	assert.Equal(t, "sit0", *ifaceList[1].Ifname)
}
