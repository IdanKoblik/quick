package networking

import (
	"context"
	"testing"
	"time"

	"quick/pkg/types"

	"github.com/go-playground/assert/v2"
)

func TestConnect(t *testing.T) {
	t.Run("rejects an unresolvable peer address", func(t *testing.T) {
		identity, err := GenerateIdentity(freeLoopbackAddr(t), types.DIRECT)
		assert.Equal(t, nil, err)

		err = Connect(context.Background(), identity, "::not-a-valid-addr::")
		assert.NotEqual(t, nil, err)
	})

	t.Run("fails to dial when no peer is listening", func(t *testing.T) {
		identity, err := GenerateIdentity(freeLoopbackAddr(t), types.DIRECT)
		assert.Equal(t, nil, err)

		// Reserved but unbound: nothing answers the handshake here.
		dead := freeLoopbackAddr(t)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		err = Connect(ctx, identity, dead)
		assert.NotEqual(t, nil, err)
	})
}
