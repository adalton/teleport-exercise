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
	"context"

	"github.com/adalton/teleport-exercise/pkg/client/jobmanager"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop a job",
	Long:    "Stop a job managed by the JobManger.  If the job is not running, this has no effect.",
	Example: "jobctl stop 8de11b74-5cd9-4769-b40d-53de13faf77f",
	RunE:    stop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stop(cmd *cobra.Command, jobIDs []string) error {

	c, err := jobmanager.NewClient(argUserID, argServerHostPort)
	if err != nil {
		return err
	}
	defer c.Close()

	for _, jobID := range jobIDs {
		err = func() error {
			ctx, cancel := context.WithTimeout(cmd.Context(), shortOperationTimeout)
			defer cancel()

			return c.Stop(ctx, jobID)
		}()

		if err != nil {
			return err
		}
	}

	return nil
}
