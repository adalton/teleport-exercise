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
	"errors"
	"os"
	"strconv"

	"github.com/adalton/teleport-exercise/pkg/client/jobmanager"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:     "query",
	Short:   "Query job state",
	Long:    "Query the state of a job managed by JobManager",
	Example: "jobctl query ba90b623-3dae-4bdd-8b96-c1ea4a999c44",
	RunE:    query,
}

func init() {
	rootCmd.AddCommand(queryCmd)
}

func query(cmd *cobra.Command, jobIDs []string) error {
	if len(jobIDs) == 0 {
		return errors.New("no jobs specified")
	}

	c, err := jobmanager.NewClient(argUserID, argServerHostPort)
	if err != nil {
		return err
	}
	defer c.Close()
	var lastError error

	jobStatusList := make([]*jobmanager.JobStatus, 0, len(jobIDs))

	for _, jobID := range jobIDs {
		// Intentionally deferring the call to cancel in the loop.  This is a
		// short-lived command and we do not expect len(jobIDs) to be large
		ctx, cancel := context.WithTimeout(cmd.Context(), shortOperationTimeout)
		defer cancel()

		status, err := c.Query(ctx, jobID)
		if err != nil {
			lastError = err
			continue
		}

		jobStatusList = append(jobStatusList, status)
	}

	renderJobStatusList(jobStatusList)

	return lastError
}

func renderJobStatusList(jobStatus []*jobmanager.JobStatus) {
	isAdmin := argUserID == jobmanager.Superuser
	header := []string{"Owner", "Name", "ID", "Running", "Pid", "Exit Code", "Signal", "Error"}

	if !isAdmin {
		header = header[1:]
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(header)

	for _, js := range jobStatus {
		runErr := ""
		if js.RunError != nil {
			runErr = js.RunError.Error()
		}

		sigStr := ""
		if js.SignalNum > 0 {
			sigStr = js.SignalNum.String()
		}

		exitCode := ""
		if js.ExitCode >= 0 {
			exitCode = strconv.FormatInt(int64(js.ExitCode), 10)
		}

		pid := ""
		if js.Pid > 0 {
			pid = strconv.FormatInt(int64(js.Pid), 10)
		}

		columns := make([]string, 0, 8)

		if isAdmin {
			columns = append(columns, js.Owner)
		}

		columns = append(columns, js.Name)
		columns = append(columns, js.ID)
		columns = append(columns, strconv.FormatBool(js.Running))
		columns = append(columns, pid)
		columns = append(columns, exitCode)
		columns = append(columns, sigStr)
		columns = append(columns, runErr)

		table.Append(columns)
	}

	table.Render()
}
