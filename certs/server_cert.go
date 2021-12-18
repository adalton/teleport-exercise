package certs

import _ "embed"

//go:embed server.cert.pem
var ServerCert []byte
