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

package certs

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"google.golang.org/grpc/credentials"
)

// NewServerTransportCredentials creates an returns a new
// credentials.TransportCredentials using the given certificate information
// with a strong TLS server configuration.
func NewServerTransportCredentials(caCert, cert, key []byte) (credentials.TransportCredentials, error) {
	const isServer = true
	return newTransportCredentials(caCert, cert, key, isServer)
}

// NewClientTransportCredentials creates an returns a new
// credentials.TransportCredentials using the given certificate information
// with a strong TLS client configuration.
func NewClientTransportCredentials(caCert, cert, key []byte) (credentials.TransportCredentials, error) {
	const isServer = false
	return newTransportCredentials(caCert, cert, key, isServer)
}

func newTransportCredentials(caCert, cert, key []byte, isServer bool) (credentials.TransportCredentials, error) {
	certificate, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("cannot append ca cert to ca pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	if isServer {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ClientCAs = capool
	} else {
		tlsConfig.RootCAs = capool
		tlsConfig.PreferServerCipherSuites = true
	}

	return credentials.NewTLS(tlsConfig), nil
}
