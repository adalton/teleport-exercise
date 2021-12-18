package certs

import _ "embed"

//go:embed badserver.cert.pem
var BadServerCert []byte
