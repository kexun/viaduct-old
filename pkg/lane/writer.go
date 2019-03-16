package lane

import "beehive/pkg/core/model"

type Writer interface {
	WriteMessage(msg *model.Message) (*model.Message, error)
}
