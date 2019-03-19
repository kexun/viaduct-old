package quic

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/kubeedge/viaduct/pkg/conn"
	quicconn "github.com/kubeedge/viaduct/pkg/conn/quic"
	"github.com/kubeedge/viaduct/pkg/mux"

	"github.com/lucas-clemente/quic-go"
)

type ConnNotify func(connection conn.Conn)

type Server struct {
	Addr               string
	TLSConfig          *tls.Config
	HandshakeTimeout   time.Duration
	ReadTimeout        time.Duration
	ReadHeaderTimeout  time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	MaxIncomingStreams int
	Handler            mux.Handler
	AutoRoute          bool
	ConnNotify         ConnNotify

	// todo
	// ConnMgr *smgr.ConnectionManager

	listener      quic.Listener
	listenerMutex sync.Mutex
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	tlsConfig, err := s.getTLSConfig(certFile, keyFile)
	if err != nil {
		return err
	}
	quicConfig := s.getQuicConfig()
	return s.serveTLS(tlsConfig, quicConfig)
}

func (s *Server) serveTLS(tlsConfig *tls.Config, quicConfig *quic.Config) error {
	s.listenerMutex.Lock()
	if s.listener != nil {
		s.listenerMutex.Unlock()
		return fmt.Errorf("ListenAndServe may only be called once")
	}

	listener, err := quic.ListenAddr(s.Addr, tlsConfig, quicConfig)
	if err != nil {
		s.listenerMutex.Unlock()
		return err
	}

	s.listener = listener
	s.listenerMutex.Unlock()

	for {
		sess, err := s.listener.Accept()
		if err != nil {
			return err
		}
		s.handleSession(sess)
	}
}

func (s *Server) handleSession(session quic.Session) {
	conn := quicconn.NewConnection(session, s.Handler)

	// todo
	//if s.ConnMgr != nil {
	//	s.ConnMgr.AddConn(conn)
	//}

	if s.ConnNotify != nil {
		go s.ConnNotify(conn)
	}

	if s.AutoRoute {
		go conn.ServeConn()
	}
}

func (s *Server) getTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	var tlsConfig *tls.Config
	if s.TLSConfig == nil {
		tlsConfig = &tls.Config{}
	} else {
		tlsConfig = s.TLSConfig.Clone()
	}

	configHasCert := len(tlsConfig.Certificates) > 0 || tlsConfig.GetCertificate != nil
	if !configHasCert || certFile != "" || keyFile != "" {
		var err error
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
	}

	return tlsConfig, nil
}

func (s *Server) getQuicConfig() *quic.Config {
	return &quic.Config{
		HandshakeTimeout:   s.HandshakeTimeout,
		IdleTimeout:        s.IdleTimeout,
		MaxIncomingStreams: s.MaxIncomingStreams,
		KeepAlive:          true,
	}
}

func (s *Server) Close() error {
	s.listenerMutex.Lock()
	defer s.listenerMutex.Unlock()

	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		return err
	}
	return nil
}
