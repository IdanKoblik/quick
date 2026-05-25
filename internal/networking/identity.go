package networking

import (
	"fmt"
	"net"
	"time"
	"crypto/rand"
	math_rand "math/rand"
	"crypto/tls"
	"encoding/pem"
	"crypto/rsa"
	"crypto/x509"
	"math/big"

	"quick/internal/utils"
	"quick/pkg/types"
)

const (
	Port = uint16(6589)
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var dictionary = []string{
	"apple", "beta", "delta", "eagle", "flame",
	"ghost", "harvest", "island", "jungle", "enter",
}

func GenerateIdentity(ip string, mode types.ConnMode) (*types.Identity, error) {
	addr, err := net.ResolveUDPAddr("udp", ip)
	if err != nil {
		return nil, err
	}

	var identity types.Identity
	identity.Addr = addr

	cert, err := GenerateCertificate()
	if err != nil {
		return nil, err
	}

	publicBytes, _ := x509.MarshalPKIXPublicKey(cert.Leaf.PublicKey)

	// The pairing code is only meaningful for P2P; direct connections are
	// established by address.
	if mode == types.P2P {
		identity.Code = GenerateCode()
	}
	identity.Fingerprint = utils.Hash(string(publicBytes))
	identity.Cert = cert

	return &identity, nil
}

func GenerateCertificate() (*tls.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048) // 2048 bits (standard)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),

		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
	}

	certDER, _ := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	return &tlsCert, nil
}

func GenerateCode() string {
	rng := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
	randomWord := dictionary[rng.Intn(len(dictionary))]

	n := 9
	ret := make([]byte, n)

	for i := range n {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}

	return fmt.Sprintf("%s-%s", randomWord, string(ret))
}
