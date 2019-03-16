package mux

import (
	"beehive/pkg/core/model"
)

type HandlerFunc func(*MessageContainer, ResponseWriter)

type ResponseWriter interface {
	WriteResponse(msg *model.Message, content interface{})
	WriteError(msg *model.Message, err string)
}

type Handler interface {
	ServeConn(msg *model.Message, writer ResponseWriter)
}

type MessageContainer struct {
	Message *model.Message
	parameters map[string]string
}

// todo
type MessageMux struct {
	moduleName string
}

var (
	MuxDefault = NewMessageMux()
)

func NewMessageMux() *MessageMux {
	return &MessageMux{}
}

func (mux *MessageMux) ServeConn(msg *model.Message, writer ResponseWriter) {
	// todo
}