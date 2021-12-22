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

package jobctl

import (
	"time"

	"github.com/spf13/cobra"
)

const (
	// On my local machine using the loopback interface, I measured the average
	// time for things like start, stop, and query to take about 190ms.
	// To allow for some transient delays and to support longer round-trips,
	// we can be conservative at ~4x that value (800ms).  If that number is too
	// conservative, we can collect more data and make a better estimate in the
	// production enviornment in which we expect clients to use the service.
	shortOperationTimeout = 800 * time.Millisecond
)

var (
	argUserID         string
	argServerHostPort string
)

var rootCmd = &cobra.Command{
	Use:   "jobctl",
	Short: "Manage JobManager jobs",
	Long:  "jobctl enables uses to manage JobManager jobs via the gRPC API",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&argUserID,
		"userID",
		"u",
		"client1",
		"The name of the user (selects the appropriate credential)")

	rootCmd.PersistentFlags().StringVar(
		&argServerHostPort,
		"hostPort",
		":24482",
		"The <hostName>:<portNumber> of the jobmanager server")
}

func Execute() error {
	return rootCmd.Execute()
}
