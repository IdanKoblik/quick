package networking

import (
	"slices"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestGenerateIdentity(t *testing.T) {
	t.Run("populates IP, port, fingerprint and code", func(t *testing.T) {
		ip := "192.168.1.42"

		identity, err := GenerateIdentity(ip)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, identity)
		assert.Equal(t, ip, identity.IP)
		assert.Equal(t, port, identity.Port)
		assert.NotEqual(t, "", identity.Fingerprint)
		assert.NotEqual(t, "", identity.Code)
	})

	t.Run("accepts empty IP string", func(t *testing.T) {
		identity, err := GenerateIdentity("")
		assert.Equal(t, nil, err)
		assert.Equal(t, "", identity.IP)
		assert.Equal(t, port, identity.Port)
	})

	t.Run("subsequent calls produce distinct fingerprints and codes", func(t *testing.T) {
		first, err := GenerateIdentity("10.0.0.1")
		assert.Equal(t, nil, err)
		second, err := GenerateIdentity("10.0.0.1")
		assert.Equal(t, nil, err)

		assert.NotEqual(t, first.Fingerprint, second.Fingerprint)
		assert.NotEqual(t, first.Code, second.Code)
	})
}

func TestGenerateFingerprint(t *testing.T) {
	t.Run("returns non-empty hex string", func(t *testing.T) {
		fp, err := GenerateFingerprint()
		assert.Equal(t, nil, err)
		assert.NotEqual(t, "", fp)
		// SHA256 hex digest is 64 characters
		assert.Equal(t, 64, len(fp))
	})

	t.Run("successive fingerprints differ", func(t *testing.T) {
		fp1, err := GenerateFingerprint()
		assert.Equal(t, nil, err)
		fp2, err := GenerateFingerprint()
		assert.Equal(t, nil, err)
		assert.NotEqual(t, fp1, fp2)
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
