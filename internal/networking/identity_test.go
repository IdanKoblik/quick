package networking

import (
	"slices"
	"strings"
	"testing"

	"quick/pkg/types"

	"github.com/go-playground/assert/v2"
)

func TestGenerateIdentity(t *testing.T) {
	t.Run("p2p populates addr, fingerprint and code", func(t *testing.T) {
		identity, err := GenerateIdentity("192.168.1.42:6589", types.P2P)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, identity)
		assert.NotEqual(t, nil, identity.Addr)
		assert.Equal(t, "192.168.1.42", identity.Addr.IP.String())
		assert.Equal(t, 6589, identity.Addr.Port)
		assert.NotEqual(t, "", identity.Fingerprint)
		assert.NotEqual(t, "", identity.Code)
	})

	t.Run("direct has no pairing code", func(t *testing.T) {
		identity, err := GenerateIdentity("192.168.1.42:6589", types.DIRECT)
		assert.Equal(t, nil, err)
		assert.Equal(t, "", identity.Code)
		assert.NotEqual(t, "", identity.Fingerprint)
	})

	t.Run("rejects address without a port", func(t *testing.T) {
		_, err := GenerateIdentity("192.168.1.42", types.DIRECT)
		assert.NotEqual(t, nil, err)
	})

	t.Run("subsequent calls produce distinct fingerprints and codes", func(t *testing.T) {
		first, err := GenerateIdentity("10.0.0.1:6589", types.P2P)
		assert.Equal(t, nil, err)
		second, err := GenerateIdentity("10.0.0.1:6589", types.P2P)
		assert.Equal(t, nil, err)

		assert.NotEqual(t, first.Fingerprint, second.Fingerprint)
		assert.NotEqual(t, first.Code, second.Code)
	})
}

func TestGenerateCode(t *testing.T) {
	t.Run("matches <word>-<9 letters> shape", func(t *testing.T) {
		code := GenerateCode()
		parts := strings.SplitN(code, "-", 2)
		assert.Equal(t, 2, len(parts))

		word, suffix := parts[0], parts[1]

		assert.Equal(t, true, slices.Contains(dictionary, word))

		assert.Equal(t, 9, len(suffix))
		for _, c := range suffix {
			assert.Equal(t, true, strings.ContainsRune(letters, c))
		}
	})

	t.Run("successive codes have distinct suffixes", func(t *testing.T) {
		c1 := GenerateCode()
		c2 := GenerateCode()
		s1 := strings.SplitN(c1, "-", 2)[1]
		s2 := strings.SplitN(c2, "-", 2)[1]
		assert.NotEqual(t, s1, s2)
	})
}
