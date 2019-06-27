package serverEventHandler

import (
	"github.com/koebeltw/Common/session"
)

type ServerEventHandler struct {
	Server session.BaseServer
	Client session.BaseClient
}

func NewServerEventHandler(BaseServer session.BaseServer, BaseClient session.BaseClient) ServerEventHandler {
	return ServerEventHandler{
		Server: BaseServer,
		Client: BaseClient,
	}
}
