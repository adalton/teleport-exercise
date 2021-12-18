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
