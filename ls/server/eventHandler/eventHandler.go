package eventHandler

import (
	"github.com/koebeltw/Common/tcp"
)

type EventHandler struct {
	Server tcp.BaseServer
	Client tcp.BaseClient
}

func NewEventHandler(BaseServer tcp.BaseServer, BaseClient tcp.BaseClient) EventHandler {
	return EventHandler{
		Server: BaseServer,
		Client: BaseClient,
	}
}
