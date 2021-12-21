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

// Note: This is a quick-n-dirty client to enable me to test the server.
//       This is _not_ intended to be the real client implementation.
//       I'll implement a more robust client in a future change.

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/adalton/teleport-exercise/certs"
	"github.com/adalton/teleport-exercise/service/jobmanager/jobmanagerv1"

	"google.golang.org/grpc"
)

const (
	// On my local machine using the loopback interface, I measured the average
	// time for things like start, stop, and query to take about 190ms.
	// To allow for some transient delays and to support longer round-trips,
	// we can be conservative at ~4x that value (800ms).  If that number is too
	// conservative, we can collect more data and make a better estimate in the
	// production enviornment in which we expect clients to use the service.
	shortOperationTimeout = 800 * time.Millisecond

	// With 1000 jobs, list takes ~200ms (~10ms more than the "short")
	//      2000 jobs, list takes ~210ms (~10ms more than 1000 jobs)
	//      3000 jobs, list takes ~220ms (~10ms more than 1000 jobs)
	// Each job adds about 0.01 ms
	// Assume a maximum of 100k jobs, that'd be ~1 second
	listOperationTimeout = shortOperationTimeout + (1 * time.Second)
)

func timeOperation(operationName string) func() {
	start := time.Now()

	return func() {
		fmt.Println(operationName, time.Since(start))
	}
}

func msgContext(timeout time.Duration) (context.Context, func()) {
	return context.WithTimeout(context.Background(), timeout)
}

func startContext() (context.Context, func()) {
	return msgContext(shortOperationTimeout)
}

func stopContext() (context.Context, func()) {
	return msgContext(shortOperationTimeout)
}

func queryContext() (context.Context, func()) {
	return msgContext(shortOperationTimeout)
}

func listContext() (context.Context, func()) {
	return msgContext(listOperationTimeout)
}

func streamContext() (context.Context, func()) {
	// There's currently no upper limit on how long the streaming operation
	// There's no need for the returned func, but keeping signature consistent
	// enables us to change our minds in the future with a simpler code change.
	return context.Background(), func() {}
}

func main() {
	user := "client1"

	if v, exists := os.LookupEnv("GRPC_USER"); exists {
		user = v
	}

	fmt.Println("User:", user)

	clientCert, clientKey, err := certs.ClientFactory(user)
	if err != nil {
		panic(err)
	}

	tc, err := certs.NewClientTransportCredentials(
		certs.CACert,
		clientCert,
		clientKey,
	)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("localhost:24482", grpc.WithTransportCredentials(tc))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)

	switch os.Args[1] {
	case "start":
		req := &jobmanagerv1.JobCreationRequest{
			Name:        os.Args[2],
			ProgramPath: os.Args[3],
			Arguments:   os.Args[4:],
		}

		ctx, cancel := startContext()
		defer cancel()

		operationDone := timeOperation("Start")
		job, err := client.Start(ctx, req)
		operationDone()

		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", job)

	case "stop":
		ctx, cancel := stopContext()
		defer cancel()

		operationDone := timeOperation("Stop")
		_, err = client.Stop(ctx, &jobmanagerv1.JobID{Id: os.Args[2]})
		operationDone()
		if err != nil {
			panic(err)
		}

	case "query", "status":
		ctx, cancel := queryContext()
		defer cancel()

		operationDone := timeOperation("Query")
		jobStatus, err := client.Query(ctx, &jobmanagerv1.JobID{Id: os.Args[2]})
		operationDone()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", jobStatus)

	case "list":
		ctx, cancel := listContext()
		defer cancel()

		operationDone := timeOperation("Query")
		statusList, err := client.List(ctx, &jobmanagerv1.NilMessage{})
		operationDone()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", statusList)

	case "stdout", "stderr":
		var streamID jobmanagerv1.OutputStream

		if os.Args[1] == "stdout" {
			streamID = jobmanagerv1.OutputStream_OutputStream_STDOUT
		} else {
			streamID = jobmanagerv1.OutputStream_OutputStream_STDERR
		}

		ctx, cancel := streamContext()
		defer cancel()

		stream, err := client.StreamOutput(ctx, &jobmanagerv1.StreamOutputRequest{
			JobID:        &jobmanagerv1.JobID{Id: os.Args[2]},
			OutputStream: streamID,
		})
		if err != nil {
			panic(err)
		}

		for {
			output, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println()
					break
				}
				panic(err)
			}

			fmt.Printf("%s", string(output.Output))
		}
	}
}
