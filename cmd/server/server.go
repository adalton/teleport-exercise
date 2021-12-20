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

package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"

	"github.com/adalton/teleport-exercise/certs"
	"github.com/adalton/teleport-exercise/pkg/command"

	"github.com/spf13/cobra"
)

var (
	argAddress string
)

var rootCmd = &cobra.Command{
	Use:   "jobmanager",
	Short: "Run the job manager server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&argAddress,
		"address",
		"a",
		":24482",
		"The <address>:<port> on which this server should listen for incoming requests")
}

func runServer() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	listener, err := net.Listen("tcp", argAddress)
	if err != nil {
		panic(err)
	}

	err = command.RunJobmanagerServer(
		ctx,
		listener,
		certs.CACert,
		certs.ServerCert,
		certs.ServerKey,
	)

	if err != nil {
		panic(err)
	}
}

func main() {
	rootCmd.Execute()
}
