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

package grpcutil

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UserIDContext is the type of the key in a context.Context where the
// userID is stored.
type UserIDContext struct{}

// UnaryGetUserIDFromContextInterceptor extracts the CommonName from the client-
// supplied certificate, creates a new context.Context with that value using
// key UserIDContext, and invokes the given handler.
func UnaryGetUserIDFromContextInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	ctx, err = getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

// streamWrapper is a wrapper over a grpc.ServerStream that enables us to
// override Context() and add the userID to it.
type streamWrapper struct {
	grpc.ServerStream
}

// Context overrides the wrapped grpc.ServerStream's Context(), finds the
// userID, and adds it to the returned context.
func (s *streamWrapper) Context() context.Context {
	ctx := s.ServerStream.Context()

	newCtx, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx
	}

	return newCtx
}

// StreamGetUserIDFromContextInterceptor is a StreamInterceptor that
// extracts the CommonName from the client-supplied certificate, creates a
// new context.Context with that value using key UserIDContext, and invokes
// the given handler.
func StreamGetUserIDFromContextInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {

	return handler(srv, &streamWrapper{ServerStream: ss})
}

// getUserIDFromContext searches the client-supplied certificates associated
// with the given ctx, extracts the CommonName, creates a new context,
// attaches the CommonName as a new value, and returns the new context.
func getUserIDFromContext(ctx context.Context) (context.Context, error) {
	if p, ok := peer.FromContext(ctx); ok {
		if mtls, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			for _, item := range mtls.State.PeerCertificates {
				if item.Subject.CommonName != "" {
					ctx = context.WithValue(ctx, &UserIDContext{}, item.Subject.CommonName)
					return ctx, nil
				}
			}
		}
	}

	return nil, status.Error(codes.Unauthenticated, "jobmanager: unauthenticated")
}
