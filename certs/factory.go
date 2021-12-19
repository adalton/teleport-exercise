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

import "fmt"

type certKeyPair struct {
	cert []byte
	key  []byte
}

var clientCertMap = map[string]certKeyPair{
	"administrator": {
		cert: AdministratorCert,
		key:  AdministratorKey,
	},
	"badclient": {
		cert: BadClientCert,
		key:  BadClientKey,
	},
	"client1": {
		cert: Client1Cert,
		key:  Client1Key,
	},
	"client2": {
		cert: Client2Cert,
		key:  Client2Key,
	},
	"client3": {
		cert: Client3Cert,
		key:  Client3Key,
	},
	"weakclient": {
		cert: WeakClientCert,
		key:  WeakClientKey,
	},
}

func ClientFactory(userID string) (cert, key []byte, err error) {
	if _, exists := clientCertMap[userID]; !exists {
		return nil, nil, fmt.Errorf("no client cert exists for userID '%s'", userID)
	}

	cert = clientCertMap[userID].cert
	key = clientCertMap[userID].key

	return
}
