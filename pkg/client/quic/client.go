package quic

import (
	"crypto/tls"
	"time"

	"github.com/kexun/viaduct/pkg/conn"
	quicconn "github.com/kexun/viaduct/pkg/conn/quic"
	"github.com/lucas-clemente/quic-go"
)

type QuicConfig struct {
	// HandshakeTimeout is the maximum duration that the cryptographic handshake may take.
	// If the timeout is exceeded, the connection is closed.
	// If this value is zero, the timeout is set to 10 seconds.
	HandshakeTimeout time.Duration
	// IdleTimeout is the maximum duration that may pass without any incoming network activity.
	// This value only applies after the handshake has completed.
	// If the timeout is exceeded, the connection is closed.
	// If this value is zero, the timeout is set to 30 seconds.
	IdleTimeout time.Duration

	// MaxIncomingStreams is the maximum number of concurrent bidirectional streams that a peer is allowed to open.
	// If not set, it will default to 100.
	// If set to a negative value, it doesn't allow any bidirectional streams.
	// Values larger than 65535 (math.MaxUint16) are invalid.
	MaxIncomingStreams int
	// MaxIncomingUniStreams is the maximum number of concurrent unidirectional streams that a peer is allowed to open.
	// This value doesn't have any effect in Google QUIC.
	// If not set, it will default to 100.
	// If set to a negative value, it doesn't allow any unidirectional streams.
	// Values larger than 65535 (math.MaxUint16) are invalid.
	MaxIncomingUniStreams int
	// KeepAlive defines whether this peer will periodically send PING frames to keep the connection alive.
	KeepAlive bool
}

type Client struct {
	Addr      string
	TLSConfig *tls.Config
	Config    *quic.Config
}

func NewClient(addr string, tlsConfig *tls.Config, quicConfig QuicConfig) *Client {
	return &Client{
		Addr:      addr,
		TLSConfig: tlsConfig,
		Config: &quic.Config{
			HandshakeTimeout:      quicConfig.HandshakeTimeout,
			IdleTimeout:           quicConfig.IdleTimeout,
			MaxIncomingStreams:    quicConfig.MaxIncomingStreams,
			MaxIncomingUniStreams: quicConfig.MaxIncomingUniStreams,
			KeepAlive:             quicConfig.KeepAlive,
		},
	}
}

func (c *Client) Dial() (conn.Conn, error) {
	session, err := quic.DialAddr(c.Addr, c.TLSConfig, c.Config)
	if err != nil {
		return nil, err
	}
	connection := quicconn.NewConnection(session, nil)
	return connection, nil
}
