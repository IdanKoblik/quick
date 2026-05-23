package networking

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"quick/internal/logging"
	"quick/pkg/types"

	"github.com/quic-go/quic-go"
)

// Connect dials a peer over QUIC using our own identity for the client
// certificate, opens a stream, and keeps the connection alive until ctx is
// cancelled.
func Connect(ctx context.Context, identity *types.Identity, peer string) error {
	addr, err := net.ResolveUDPAddr("udp4", peer)
	if err != nil {
		return fmt.Errorf("resolve peer %q: %w", peer, err)
	}

	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return err
	}
	defer udpConn.Close()

	tr := &quic.Transport{Conn: udpConn}
	defer tr.Close()

	tlsConf := &tls.Config{
		Certificates:       []tls.Certificate{*identity.Cert},
		NextProtos:         []string{"p2p-file-transfer"},
		InsecureSkipVerify: true,
	}

	conn, err := tr.Dial(ctx, addr, tlsConf, quicConf)
	if err != nil {
		return fmt.Errorf("dial peer %s: %w", addr, err)
	}
	defer conn.CloseWithError(0, "")

	remote := conn.RemoteAddr()
	logging.Log.Infof("connected to peer: %s", remote)
	defer logging.Log.Infof("disconnected from peer: %s", remote)

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return fmt.Errorf("open stream to %s: %w", remote, err)
	}
	defer stream.Close()

	logging.Log.Debugf("stream %d opened to %s", stream.StreamID(), remote)

	select {
	case <-ctx.Done():
		return nil
	case <-conn.Context().Done():
		return nil
	}
}
