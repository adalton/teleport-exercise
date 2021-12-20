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

package server_test

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/adalton/teleport-exercise/certs"
	"github.com/adalton/teleport-exercise/pkg/command"
	"github.com/adalton/teleport-exercise/service/jobmanager/jobmanagerv1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_clientServer_clientCertNotSignedByTrustedCA(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	hostPort, err := runServer(ctx, &wg, t, certs.CACert, certs.ServerCert, certs.ServerKey)
	require.Nil(t, err)

	tc, err := certs.NewClientTransportCredentials(
		certs.CACert,
		certs.BadClientCert,
		certs.BadClientKey,
	)
	require.Nil(t, err)

	conn, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(tc))
	require.Nil(t, err)
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)
	_, err = client.List(context.Background(), &jobmanagerv1.NilMessage{})

	assert.Error(t, err)
	s, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unavailable, s.Code())

	cancel()
	wg.Wait()
}

func Test_clientServer_serverCertNotSignedByTrustedCA(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	hostPort, err := runServer(ctx, &wg, t, certs.CACert, certs.BadServerCert, certs.BadServerKey)
	require.Nil(t, err)

	tc, err := certs.NewClientTransportCredentials(
		certs.CACert,
		certs.Client1Cert,
		certs.Client1Key,
	)
	require.Nil(t, err)

	conn, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(tc))
	require.Nil(t, err)
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)
	_, err = client.List(context.Background(), &jobmanagerv1.NilMessage{})

	assert.Error(t, err)
	s, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unavailable, s.Code())

	cancel()
	wg.Wait()
}

func Test_clientServer_TooWeakServerCert(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	hostPort, err := runServer(ctx, &wg, t, certs.CACert, certs.WeakServerCert, certs.WeakServerKey)
	require.Nil(t, err)

	tc, err := certs.NewClientTransportCredentials(
		certs.CACert,
		certs.Client1Cert,
		certs.Client1Key,
	)
	require.Nil(t, err)

	conn, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(tc))
	require.Nil(t, err)
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)
	_, err = client.List(context.Background(), &jobmanagerv1.NilMessage{})

	assert.Error(t, err)

	s, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unavailable, s.Code())

	cancel()
	wg.Wait()
}

func Test_clientServer_TooWeakClientCert(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	hostPort, err := runServer(ctx, &wg, t, certs.CACert, certs.ServerCert, certs.ServerKey)
	require.Nil(t, err)

	tc, err := certs.NewClientTransportCredentials(
		certs.CACert,
		certs.WeakClientCert,
		certs.WeakClientKey,
	)
	require.Nil(t, err)

	conn, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(tc))
	require.Nil(t, err)
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)
	_, err = client.List(context.Background(), &jobmanagerv1.NilMessage{})

	assert.Error(t, err)
	s, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unavailable, s.Code())

	cancel()
	wg.Wait()
}

func Test_clientServer_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	hostPort, err := runServer(ctx, &wg, t, certs.CACert, certs.ServerCert, certs.ServerKey)
	require.Nil(t, err)

	tc, err := certs.NewClientTransportCredentials(
		certs.CACert,
		certs.Client1Cert,
		certs.Client1Key,
	)
	require.Nil(t, err)

	conn, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(tc))
	require.Nil(t, err)
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)
	_, err = client.List(context.Background(), &jobmanagerv1.NilMessage{})

	assert.Nil(t, err)

	cancel()
	wg.Wait()
}

func Test_clientServer_Multitenant(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	hostPort, err := runServer(ctx, &wg, t, certs.CACert, certs.ServerCert, certs.ServerKey)
	require.Nil(t, err)

	tcClient1, err := certs.NewClientTransportCredentials(
		certs.CACert,
		certs.Client1Cert,
		certs.Client1Key,
	)
	require.Nil(t, err)

	connClient1, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(tcClient1))
	require.Nil(t, err)
	defer connClient1.Close()

	client := jobmanagerv1.NewJobManagerClient(connClient1)
	_, err = client.Start(context.Background(), &jobmanagerv1.JobCreationRequest{
		Name:        "myjob",
		ProgramPath: "/bin/true",
	})
	assert.Nil(t, err)

	tcClient2, err := certs.NewClientTransportCredentials(
		certs.CACert,
		certs.Client2Cert,
		certs.Client2Key,
	)
	require.Nil(t, err)

	connClient2, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(tcClient2))
	require.Nil(t, err)
	defer connClient2.Close()

	client2 := jobmanagerv1.NewJobManagerClient(connClient2)
	jobList, err := client2.List(context.Background(), &jobmanagerv1.NilMessage{})

	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobList.JobStatusList))

	cancel()
	wg.Wait()
}

func Test_clientServer_AdministratorCanSeeAllJobs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	hostPort, err := runServer(ctx, &wg, t, certs.CACert, certs.ServerCert, certs.ServerKey)
	require.Nil(t, err)

	tcClient1, err := certs.NewClientTransportCredentials(
		certs.CACert,
		certs.Client1Cert,
		certs.Client1Key,
	)
	require.Nil(t, err)

	connClient1, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(tcClient1))
	require.Nil(t, err)
	defer connClient1.Close()

	client := jobmanagerv1.NewJobManagerClient(connClient1)
	_, err = client.Start(context.Background(), &jobmanagerv1.JobCreationRequest{
		Name:        "myjob",
		ProgramPath: "/bin/true",
	})
	assert.Nil(t, err)

	tcClient2, err := certs.NewClientTransportCredentials(
		certs.CACert,
		certs.AdministratorCert,
		certs.AdministratorKey,
	)
	require.Nil(t, err)

	connClient2, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(tcClient2))
	require.Nil(t, err)
	defer connClient2.Close()

	client2 := jobmanagerv1.NewJobManagerClient(connClient2)
	jobList, err := client2.List(context.Background(), &jobmanagerv1.NilMessage{})

	assert.Nil(t, err)
	assert.Equal(t, 1, len(jobList.JobStatusList))

	cancel()
	wg.Wait()
}

func runServer(
	ctx context.Context,
	wg *sync.WaitGroup,
	t *testing.T,
	caCert, serverCert, serverKey []byte,
) (hostPort string, err error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	go func() {
		runErr := command.RunJobmanagerServer(
			ctx, listener, caCert, serverCert, serverKey)
		if runErr != nil {
			t.Error()
		}
		wg.Done()
	}()

	return listener.Addr().String(), nil
}
