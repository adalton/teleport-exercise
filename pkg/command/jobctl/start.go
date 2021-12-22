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
	"os"

	"github.com/adalton/teleport-exercise/pkg/client/jobmanager"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	argStartJobName string
	argJobCommand   string
)

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a new job",
	Long:    "Starts a new job with the given parameters on the JobManager",
	Example: "start -j myJob -c /usr/bin/find -- /dir -type f",
	RunE:    start,
}

func init() {
	startCmd.PersistentFlags().StringVarP(
		&argStartJobName,
		"jobName",
		"j",
		"",
		"The name of the job to create; must be unique",
	)
	startCmd.MarkPersistentFlagRequired("jobName")

	startCmd.PersistentFlags().StringVarP(
		&argJobCommand,
		"command",
		"c",
		"",
		"The command for the job to run; must supply full path",
	)
	startCmd.MarkPersistentFlagRequired("command")

	rootCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(cmd.Context(), shortOperationTimeout)
	defer cancel()

	c, err := jobmanager.NewClient(argUserID, argServerHostPort)
	if err != nil {
		return err
	}
	defer c.Close()

	jobID, err := c.Start(ctx, argStartJobName, argJobCommand, args...)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Name", "ID"})
	table.Append([]string{argStartJobName, jobID})

	table.Render()

	return nil
}
