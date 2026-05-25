package types

import (
	"crypto/tls"
	"net"
)

type Identity struct {
	Addr 		*net.UDPAddr
	Fingerprint string
	Code 		string
	Cert 		*tls.Certificate
}
