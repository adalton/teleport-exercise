package certs

import _ "embed"

//go:embed weakclient.key.pem
var WeakClientKey []byte
