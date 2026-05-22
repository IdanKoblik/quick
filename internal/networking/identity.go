package networking

import (
	"fmt"
	"time"
	"crypto/rand"
	math_rand "math/rand"
	"crypto/rsa"
	"crypto/x509"
	"math/big"

	"quick/internal/utils"
	"quick/pkg/types"
)

const (
	port = uint16(6589)
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var dictionary = []string{
	"apple", "beta", "delta", "eagle", "flame",
	"ghost", "harvest", "island", "jungle", "enter",
}

func GenerateIdentity(ip string) (*types.Identity, error) {
	var identity types.Identity
	identity.IP = ip
	identity.Port = port

	fp, err := GenerateFingerprint()
	if err != nil {
		return nil, err
	}

	code := GenerateCode()
	identity.Code = code
	identity.Fingerprint = fp

	return &identity, nil
}

func GenerateFingerprint() (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048) // 2048 bits (standard)
	if err != nil {
		return "", err
	}

	publicKey := &privateKey.PublicKey
	publicBytes, _ := x509.MarshalPKIXPublicKey(publicKey)

	return utils.Hash(string(publicBytes)), nil
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
