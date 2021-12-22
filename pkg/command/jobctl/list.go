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
	"fmt"
	"time"

	"github.com/adalton/teleport-exercise/pkg/client/jobmanager"
	"github.com/spf13/cobra"
)

const (
	// With 1000 jobs, list takes ~200ms (~10ms more than the "short")
	//      2000 jobs, list takes ~210ms (~10ms more than 1000 jobs)
	//      3000 jobs, list takes ~220ms (~10ms more than 1000 jobs)
	// Each job adds about 0.01 ms
	// Assume a maximum of 100k jobs, that'd be ~1 second
	listOperationTimeout = shortOperationTimeout + (1 * time.Second)
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List job",
	Long:    "List jobs managed by the JobManager",
	Example: "jobctl list",
	RunE:    list,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(cmd.Context(), listOperationTimeout)
	defer cancel()

	c, err := jobmanager.NewClient(argUserID, argServerHostPort)
	if err != nil {
		return err
	}
	defer c.Close()

	jobList, err := c.List(ctx)
	if err != nil {
		return err
	}

	if len(jobList) == 0 {
		fmt.Println("There are no jobs")
		return nil
	}

	renderJobStatusList(jobList)

	return nil
}
