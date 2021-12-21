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

package jobmanager

import (
	"context"
	"errors"
	"io"
	"syscall"

	"github.com/adalton/teleport-exercise/certs"
	"github.com/adalton/teleport-exercise/pkg/jobmanager"
	"github.com/adalton/teleport-exercise/service/jobmanager/jobmanagerv1"

	"google.golang.org/grpc"
)

// JobStatus models the current status of a job.
type JobStatus = jobmanager.JobStatus

// Superuser is the name of the user who can access any job.
const Superuser = jobmanager.Superuser

// Client is a client interface to the JobManager gRPC server.  It isolates
// code that want to communicate with the JobManager server from the gRPC details.
type Client struct {
	jm   jobmanagerv1.JobManagerClient
	conn *grpc.ClientConn
}

// NewClient creates a new JobManager client using the given user's credentials
// and connecting to a JobManager server at the given hostPort.
func NewClient(userID, hostPort string) (*Client, error) {
	clientCert, clientKey, err := certs.ClientFactory(userID)
	if err != nil {
		return nil, err
	}

	tc, err := certs.NewClientTransportCredentials(certs.CACert, clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(tc))
	if err != nil {
		return nil, err
	}

	return &Client{
		jm:   jobmanagerv1.NewJobManagerClient(conn),
		conn: conn,
	}, nil
}

// Start invokes an RPC on the JobManager server to start a new job.
func (c *Client) Start(
	ctx context.Context,
	jobName, programPath string,
	programArgs ...string,
) (jobID string, err error) {

	job, err := c.jm.Start(ctx, &jobmanagerv1.JobCreationRequest{
		Name:        jobName,
		ProgramPath: programPath,
		Arguments:   programArgs,
	})
	if err != nil {
		return "", err
	}

	return job.Id.Id, nil
}

// Stop invokes an RPC on the JobManager to stop a job.  If the job with the
// given jobID isn't running, this operation does nothing.
func (c *Client) Stop(ctx context.Context, jobID string) error {
	_, err := c.jm.Stop(ctx, &jobmanagerv1.JobID{Id: jobID})

	return err
}

// Query invokes an RPC on the JobManager server to retrieve the current status
// of the job with the given jobID.
func (c *Client) Query(ctx context.Context, jobID string) (*JobStatus, error) {
	jobStatus, err := c.jm.Query(ctx, &jobmanagerv1.JobID{Id: jobID})
	if err != nil {
		return nil, err
	}

	return jobStatusRpcToLocal(jobStatus), nil
}

// Query invokes an RPC on the JobManager server to retrieve the list of jobs
// started by the user.  If the user is the administrator, then it returns a
// list of all jobs in the system.
func (c *Client) List(ctx context.Context) ([]*JobStatus, error) {
	jobStatusList, err := c.jm.List(ctx, &jobmanagerv1.NilMessage{})
	if err != nil {
		return nil, err
	}

	retList := make([]*JobStatus, len(jobStatusList.JobStatusList))

	for i, js := range jobStatusList.JobStatusList {
		retList[i] = jobStatusRpcToLocal(js)
	}

	return retList, nil
}

// StreamStdout invokes an RPC on the JobManager server to stream the standard
// output of the job with the given jobID.  This function will block until either
// (1) the context is interrupted, or (2) the job completes.
func (c *Client) StreamStdout(ctx context.Context, jobID string, out io.Writer) error {
	return c.stream(ctx, jobID, out, jobmanagerv1.OutputStream_OutputStream_STDOUT)
}

// StreamStdout invokes an RPC on the JobManager server to stream the standard
// error of the job with the given jobID.  This function will block until either
// (1) the context is interrupted, or (2) the job completes.
func (c *Client) StreamStderr(ctx context.Context, jobID string, out io.Writer) error {
	return c.stream(ctx, jobID, out, jobmanagerv1.OutputStream_OutputStream_STDERR)
}

// Close closes the connection to the JobManager server.
func (c *Client) Close() {
	c.conn.Close()
}

func jobStatusRpcToLocal(jobStatus *jobmanagerv1.JobStatus) *JobStatus {
	var runError error
	if jobStatus.ErrorMessage != "" {
		runError = errors.New(jobStatus.ErrorMessage)
	}

	return &JobStatus{
		Owner:     jobStatus.Owner,
		Name:      jobStatus.Job.Name,
		ID:        jobStatus.Job.Id.Id,
		Running:   jobStatus.IsRunning,
		Pid:       int(jobStatus.Pid),
		ExitCode:  int(jobStatus.ExitCode),
		SignalNum: syscall.Signal(jobStatus.SignalNumber),
		RunError:  runError,
	}
}

func (c *Client) stream(
	ctx context.Context,
	jobID string,
	out io.Writer,
	outputStream jobmanagerv1.OutputStream,
) error {

	grpcStream, err := c.jm.StreamOutput(ctx, &jobmanagerv1.StreamOutputRequest{
		JobID:        &jobmanagerv1.JobID{Id: jobID},
		OutputStream: outputStream,
	})
	if err != nil {
		return err
	}

	for {
		output, err := grpcStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		out.Write(output.Output)
	}

	return nil
}
