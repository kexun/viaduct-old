package quiclane

import (
	"time"

	"github.com/kubeedge/viaduct/pkg/lane"
	"github.com/lucas-clemente/quic-go"

	"beehive/pkg/core/model"
)

const NonDelay = 0

type writer struct {
	stream quic.Stream
	timeout time.Duration
}

func NewQuicWriter(stream quic.Stream) lane.Writer {
	return &writer{
		stream: stream,
	}
}

func NewQuicWriterTimeout(stream quic.Stream, timeout time.Duration) lane.Writer {
	return &writer{
		stream: stream,
		timeout: timeout,
	}
}

func (w *writer) WriteMessage(msg *model.Message) (*model.Message, error) {

	// todo
	return nil, nil
}

func (w *writer) write(data []byte) error {

	// todo
	return nil
}