package quiclane

import (
	"fmt"

	"beehive/pkg/core/model"
	"github.com/kubeedge/viaduct/pkg/lane"

	"github.com/lucas-clemente/quic-go"
)

type reader struct {
	stream        quic.Stream
	payloadBuffer []byte
	payloadMax    uint32
}

const payloadMaxSize = uint32(1024 * 1024)

func NewQuicReader(stream quic.Stream) lane.Reader {
	return &reader{
		stream:     stream,
		payloadMax: payloadMaxSize,
	}
}

func (r *reader) ReadMessage(msg *model.Message) error {

	// todo
	return nil
}

func (r *reader) read() ([]byte, error) {
	if r.stream == nil {
		return nil, fmt.Errorf("quic reader stream is nil")
	}

	// todo
	return nil, nil
}
