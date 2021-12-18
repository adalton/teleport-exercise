package certs

import _ "embed"

//go:embed ca.cert.pem
var CACert []byte
