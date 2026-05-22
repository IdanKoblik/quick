package networking

import (
	"net"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestGetInterfaces(t *testing.T) {
	t.Run("nil pointer returns empty map", func(t *testing.T) {
		result := GetInterfaces(nil)
		assert.NotEqual(t, nil, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("empty slice returns empty map", func(t *testing.T) {
		ifaces := []net.Interface{}
		result := GetInterfaces(&ifaces)
		assert.Equal(t, 0, len(result))
	})

	t.Run("returns addresses for real system interfaces", func(t *testing.T) {
		ifaces, err := net.Interfaces()
		assert.Equal(t, nil, err)
		if len(ifaces) == 0 {
			t.Skip("no network interfaces available on host")
		}

		result := GetInterfaces(&ifaces)

		// Every interface that reported at least one address must be present
		// in the map with the same address strings, in order.
		for _, iface := range ifaces {
			addrs, _ := iface.Addrs()
			if len(addrs) == 0 {
				continue
			}

			got, ok := result[iface.Name]
			assert.Equal(t, true, ok)
			assert.Equal(t, len(addrs), len(got))
			for i, addr := range addrs {
				assert.Equal(t, addr.String(), got[i])
			}
		}
	})
}
