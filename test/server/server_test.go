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
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/adalton/teleport-exercise/pkg/command"
	"github.com/adalton/teleport-exercise/service/jobmanager/jobmanagerv1"
	"github.com/adalton/teleport-exercise/util/grpcutil"
	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: There is some fragility with these tests.
//       * This expects port 12345 to be available (so we wouldn't want to
//         run multiple instances of this test on the same host in parallel)
//       * There's a Sleep() between starting the server and trying to connect
//         to it.  That's open to race conditions.  I'd like to find a way to
//         know that the server is up before we try to connect to it.

var certDir string

func init() {
	dir, ok := os.LookupEnv("CERTDIR")
	if ok {
		certDir = dir
	}
}

func Test_clientServer_clientCertNotSignedByTrustedCA(t *testing.T) {
	if certDir == "" {
		t.Skip("Skipping test, CERTDIR not set")
	}

	stop := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		_ = command.RunJobmanagerServer(
			"tcp",
			":12345",
			certDir+"/ca.cert.pem",
			certDir+"/server.cert.pem",
			certDir+"/server.key.pem",
			stop)
		wg.Done()
	}()

	waitForServer()

	tc, err := grpcutil.NewClientTransportCredentials(
		certDir+"/ca.cert.pem",
		certDir+"/badclient.cert.pem",
		certDir+"/badclient.key.pem",
	)
	require.Nil(t, err)

	conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(tc))
	require.Nil(t, err)
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)

	ctx := grpcutil.AttachUserIDToContext(context.Background(), "user1")

	_, err = client.List(ctx, &jobmanagerv1.NilMessage{})

	assert.Error(t, err)

	stop <- os.Kill
	wg.Wait()
}

func Test_clientServer_serverCertNotSignedByTrustedCA(t *testing.T) {
	if certDir == "" {
		t.Skip("Skipping test, CERTDIR not set")
	}

	stop := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		_ = command.RunJobmanagerServer(
			"tcp",
			":12345",
			certDir+"/ca.cert.pem",
			certDir+"/badserver.cert.pem",
			certDir+"/badserver.key.pem",
			stop)
		wg.Done()
	}()

	waitForServer()

	tc, err := grpcutil.NewClientTransportCredentials(
		certDir+"/ca.cert.pem",
		certDir+"/client1.cert.pem",
		certDir+"/client1.key.pem",
	)
	require.Nil(t, err)

	conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(tc))
	require.Nil(t, err)
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)

	ctx := grpcutil.AttachUserIDToContext(context.Background(), "user1")

	_, err = client.List(ctx, &jobmanagerv1.NilMessage{})

	fmt.Println(err)
	assert.Error(t, err)

	stop <- os.Kill
	wg.Wait()
}

func Test_clientServer_TooWeakServerCert(t *testing.T) {
	if certDir == "" {
		t.Skip("Skipping test, CERTDIR not set")
	}

	stop := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		_ = command.RunJobmanagerServer(
			"tcp",
			":12345",
			certDir+"/ca.cert.pem",
			certDir+"/weakserver.cert.pem",
			certDir+"/weakserver.key.pem",
			stop)
		wg.Done()
	}()

	waitForServer()

	tc, err := grpcutil.NewClientTransportCredentials(
		certDir+"/ca.cert.pem",
		certDir+"/client1.cert.pem",
		certDir+"/client1.key.pem",
	)
	require.Nil(t, err)

	conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(tc))
	require.Nil(t, err)
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)

	ctx := grpcutil.AttachUserIDToContext(context.Background(), "user1")

	_, err = client.List(ctx, &jobmanagerv1.NilMessage{})

	fmt.Println(err)
	assert.Error(t, err)

	stop <- os.Kill
	wg.Wait()
}

func Test_clientServer_TooWeakClientCert(t *testing.T) {
	if certDir == "" {
		t.Skip("Skipping test, CERTDIR not set")
	}

	stop := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		_ = command.RunJobmanagerServer(
			"tcp",
			":12345",
			certDir+"/ca.cert.pem",
			certDir+"/server.cert.pem",
			certDir+"/server.key.pem",
			stop)
		wg.Done()
	}()

	waitForServer()

	tc, err := grpcutil.NewClientTransportCredentials(
		certDir+"/ca.cert.pem",
		certDir+"/weakclient.cert.pem",
		certDir+"/weakclient.key.pem",
	)
	require.Nil(t, err)

	conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(tc))
	require.Nil(t, err)
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)

	ctx := grpcutil.AttachUserIDToContext(context.Background(), "weakclient")

	_, err = client.List(ctx, &jobmanagerv1.NilMessage{})

	fmt.Println(err)
	assert.Error(t, err)

	stop <- os.Kill
	wg.Wait()
}

func Test_clientServer_Success(t *testing.T) {
	if certDir == "" {
		t.Skip("Skipping test, CERTDIR not set")
	}

	stop := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		_ = command.RunJobmanagerServer(
			"tcp",
			":12345",
			certDir+"/ca.cert.pem",
			certDir+"/server.cert.pem",
			certDir+"/server.key.pem",
			stop)
		wg.Done()
	}()

	waitForServer()

	tc, err := grpcutil.NewClientTransportCredentials(
		certDir+"/ca.cert.pem",
		certDir+"/client1.cert.pem",
		certDir+"/client1.key.pem",
	)
	require.Nil(t, err)

	conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(tc))
	require.Nil(t, err)
	defer conn.Close()

	client := jobmanagerv1.NewJobManagerClient(conn)

	ctx := grpcutil.AttachUserIDToContext(context.Background(), "user1")

	_, err = client.List(ctx, &jobmanagerv1.NilMessage{})

	assert.Nil(t, err)

	stop <- os.Kill
	wg.Wait()
}

// waitForServer waits for the server to come up before attempting a network
// connection to it.
//
// This is gross, but I didn't find a better way to check.
func waitForServer() {
	time.Sleep(1 * time.Second)
}
