package lane

import "beehive/pkg/core/model"

type Reader interface {
	ReadMessage(msg *model.Message) error
}
