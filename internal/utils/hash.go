package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(raw string) string {
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}
