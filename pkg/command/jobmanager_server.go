/*
 *
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

package command

import (
	"context"
	"log"
	"net"

	"github.com/adalton/teleport-exercise/certs"
	"github.com/adalton/teleport-exercise/server/jobmanager/serverv1"
	"github.com/adalton/teleport-exercise/service/jobmanager/jobmanagerv1"

	"google.golang.org/grpc"
)

// RunJobmanagerServer runs a JobmanagerServer on the given network, listening
// on the given address, with the given CA certificate and server certificate
// and key.
func RunJobmanagerServer(
	ctx context.Context,
	listener net.Listener,
	caCert, serverCert, serverKey []byte,
) error {
	tc, err := certs.NewServerTransportCredentials(caCert, serverCert, serverKey)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tc),
		grpc.UnaryInterceptor(serverv1.UnaryGetUserIDFromContextInterceptor),
		grpc.StreamInterceptor(serverv1.StreamGetUserIDFromContextInterceptor),
	)

	jobmanagerv1.RegisterJobManagerServer(grpcServer, serverv1.NewJobmanagerServer())

	errChan := make(chan error)

	go func() {
		log.Println("Server ready.")
		if err := grpcServer.Serve(listener); err != nil {
			errChan <- err
			return
		}
	}()

	select {
	case <-ctx.Done():
	case err = <-errChan:
	}

	grpcServer.GracefulStop()

	return err
}
