package networking

import (
	"context"
	"net"
	"testing"
	"time"

	"quick/pkg/types"

	"github.com/go-playground/assert/v2"
)

// freeLoopbackAddr reserves an ephemeral UDP port on the loopback interface and
// returns its "host:port". The socket is closed before returning so the caller
// can bind it; the window is small enough to be reliable for tests.
func freeLoopbackAddr(t *testing.T) string {
	t.Helper()

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatalf("reserve udp port: %v", err)
	}
	addr := conn.LocalAddr().String()
	conn.Close()

	return addr
}

func TestStartServer(t *testing.T) {
	t.Run("returns nil when the context is cancelled", func(t *testing.T) {
		identity, err := GenerateIdentity(freeLoopbackAddr(t), types.DIRECT)
		assert.Equal(t, nil, err)

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error, 1)
		go func() { done <- StartServer(ctx, identity) }()

		// Let the listener bind, then ask it to shut down.
		time.Sleep(150 * time.Millisecond)
		cancel()

		select {
		case err := <-done:
			assert.Equal(t, nil, err)
		case <-time.After(2 * time.Second):
			t.Fatal("StartServer did not return after context cancellation")
		}
	})

	t.Run("fails when the address is already in use", func(t *testing.T) {
		addr := freeLoopbackAddr(t)

		// Hold the port so the server cannot bind it.
		udpAddr, err := net.ResolveUDPAddr("udp4", addr)
		assert.Equal(t, nil, err)
		held, err := net.ListenUDP("udp4", udpAddr)
		assert.Equal(t, nil, err)
		defer held.Close()

		identity, err := GenerateIdentity(addr, types.DIRECT)
		assert.Equal(t, nil, err)

		err = StartServer(context.Background(), identity)
		assert.NotEqual(t, nil, err)
	})
}

func TestServerClientRoundTrip(t *testing.T) {
	addr := freeLoopbackAddr(t)

	server, err := GenerateIdentity(addr, types.DIRECT)
	assert.Equal(t, nil, err)

	serverCtx, stopServer := context.WithCancel(context.Background())
	defer stopServer()

	serverDone := make(chan error, 1)
	go func() { serverDone <- StartServer(serverCtx, server) }()

	// Let the listener bind before dialing.
	time.Sleep(150 * time.Millisecond)

	client, err := GenerateIdentity(freeLoopbackAddr(t), types.DIRECT)
	assert.Equal(t, nil, err)

	clientCtx, disconnect := context.WithCancel(context.Background())
	clientDone := make(chan error, 1)
	go func() { clientDone <- Connect(clientCtx, client, addr) }()

	// Once connected, Connect blocks until its context is cancelled. Give the
	// handshake time to complete, then hang up: a successful round trip returns nil.
	time.Sleep(500 * time.Millisecond)
	disconnect()

	select {
	case err := <-clientDone:
		assert.Equal(t, nil, err)
	case <-time.After(3 * time.Second):
		t.Fatal("client did not disconnect")
	}

	stopServer()
	select {
	case err := <-serverDone:
		assert.Equal(t, nil, err)
	case <-time.After(3 * time.Second):
		t.Fatal("server did not stop")
	}
}
