package conn

import (
	"time"
	
	"beehive/pkg/core/model"
)

type Conn interface {
	ServeConn()
	WriteMessageSync(message *model.Message, timeout time.Duration) (*model.Message, error)
	WriteMessageAsyn(message *model.Message) (*model.Message, error)
	ReadMessage(message *model.Message) error
	Close() error
	CloseWithError(code uint16, err error) error
}
