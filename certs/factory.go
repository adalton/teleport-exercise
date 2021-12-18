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

/*
var serverCertMap = map[string]certKeyPair{
	"badserver": {
		cert: BadServerCert,
		key:  BadServerKey,
	},
	"server": {
		cert: ServerCert,
		key:  ServerKey,
	},
	"weakserver": {
		cert: WeakServerCert,
		key:  WeakServerKey,
	},
}
*/

func ClientFactory(userID string) (cert, key []byte, err error) {
	if _, exists := clientCertMap[userID]; !exists {
		return nil, nil, fmt.Errorf("no client cert exists for userID '%s'", userID)
	}

	cert = clientCertMap[userID].cert
	key = clientCertMap[userID].key

	return
}

/*
func ServertCertFactory(name string) ([]byte, error) {
	if _, exists := serverCertMap[name]; !exists {
		return nil, fmt.Errorf("no server cert exists for name '%s'", name)
	}

	return serverCertMap[name].cert, nil
}

func ServertKeyFactory(name string) ([]byte, error) {
	if _, exists := serverCertMap[name]; !exists {
		return nil, fmt.Errorf("no server key exists for name '%s'", name)
	}

	return serverCertMap[name].key, nil
}
*/

/*
badca.cert.pem
ca.cert.pem
*/
