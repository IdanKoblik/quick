package networking

import (
	"net"
	"time"
	"context"
	"crypto/tls"

	"quick/internal/logging"
	"quick/pkg/types"

	"github.com/quic-go/quic-go"
)

var quicConf = &quic.Config{
	MaxIdleTimeout: 30 * time.Second,
	KeepAlivePeriod: 10 * time.Second,
}

func StartServer(ctx context.Context, identity *types.Identity) error {
	udpConn, err := net.ListenUDP("udp4", identity.Addr)
	if err != nil {
		return err
	}
	defer udpConn.Close()

	tr := &quic.Transport{
		Conn: udpConn,
	}
	defer tr.Close()

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{*identity.Cert},
		NextProtos:   []string{"p2p-file-transfer"},
	}

	listener, err := tr.Listen(tlsConf, quicConf)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			logging.Log.Error(err)
			continue
		}

		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn *quic.Conn) {
	defer conn.CloseWithError(0, "")

	remote := conn.RemoteAddr()
	logging.Log.Infof("peer connected: %s", remote)
	defer logging.Log.Infof("peer disconnected: %s", remote)

	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			if ctx.Err() != nil || conn.Context().Err() != nil {
				return
			}

			logging.Log.Errorf("accept stream from %s: %v", remote, err)
			return
		}

		go handleStream(stream, remote)
	}
}

func handleStream(stream *quic.Stream, remote net.Addr) {
	defer stream.Close()
	logging.Log.Debugf("stream %d opened from %s", stream.StreamID(), remote)
}

