package quic

import (
	"fmt"
	"time"

	"beehive/pkg/common/log"
	"beehive/pkg/core/model"
	"github.com/kubeedge/viaduct/pkg/conn"
	"github.com/kubeedge/viaduct/pkg/lane/quiclane"
	"github.com/kubeedge/viaduct/pkg/mux"

	"github.com/lucas-clemente/quic-go"
)

type connection struct {
	Session quic.Session
	Handler mux.Handler
}

func NewConnection(session quic.Session, handler mux.Handler) conn.Conn {
	return &connection{
		Session: session,
		Handler: handler,
	}
}

func (conn *connection) ServeConn() {
	for {
		stream, err := conn.Session.AcceptStream()
		if err != nil {
			conn.Session.Close()
			return
		}
		go conn.handleMessage(stream)
	}
}

func (conn *connection) handleMessage(stream quic.Stream) {
	msg := &model.Message{}
	for {
		err := quiclane.NewQuicReader(stream).ReadMessage(msg)
		if err != nil {
			log.LOGGER.Errorf("Failed to read message %s", err)
			stream.Close()
			return
		}
		if conn.Handler == nil {
			conn.Handler = mux.MuxDefault
		}
		conn.Handler.ServeConn(msg, &responseWriter{stream})
	}
}

func (conn *connection) Close() error {
	return conn.Session.Close()
}

func (conn *connection) CloseWithError(code uint16, err error) error {
	return conn.Session.CloseWithError(quic.ErrorCode(code), err)
}

func (conn *connection) WriteMessageSync(message *model.Message, timeout time.Duration) (*model.Message, error) {
	message.Header.Sync = true
	return conn.writeMessage(message, timeout)
}

func (conn *connection) WriteMessageAsyn(message *model.Message) (*model.Message, error) {
	message.Header.Sync = false
	return conn.writeMessage(message, quiclane.NonDelay)
}

func (conn *connection) writeMessage(message *model.Message, timeout time.Duration) (*model.Message, error) {
	if conn.Session == nil {
		return nil, fmt.Errorf("bad conn session")
	}

	stream, err := conn.Session.OpenStreamSync()
	if err != nil {
		conn.Session.Close()
		return nil, fmt.Errorf("failed to open stream")
	}
	defer stream.Close()

	return quiclane.NewQuicWriterTimeout(stream, timeout).WriteMessage(message)
}

func (conn *connection) ReadMessage(message *model.Message) error {
	if conn.Session == nil {
		return fmt.Errorf("bad conn session")
	}

	stream, err := conn.Session.AcceptStream()
	if err != nil {
		log.LOGGER.Errorf("ailed to accept stream [%v]", err)
		return err
	}

	return quiclane.NewQuicReader(stream).ReadMessage(message)
}

type responseWriter struct{ quic.Stream }

func (r *responseWriter) WriteResponse(msg *model.Message, content interface{}) {
	response := msg.NewRespByMessage(msg, content)
	response.Header.Sync = false
	_, err := quiclane.NewQuicWriter(r.Stream).WriteMessage(response)
	if err != nil {
		log.LOGGER.Errorf("Failed to write response [%v]", err)
	}
}

func (r *responseWriter) WriteError(msg *model.Message, errMsg string) {
	response := model.NewErrorMessage(msg, errMsg)
	response.Header.Sync = false
	_, err := quiclane.NewQuicWriter(r.Stream).WriteMessage(response)
	if err != nil {
		log.LOGGER.Errorf("Failed to write err response [%v]", err)
	}
}
