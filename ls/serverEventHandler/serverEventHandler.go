package serverEventHandler

import (
	"github.com/koebeltw/Common/session"
)

type ServerEventHandler struct {
	BaseServer session.BaseServer
	BaseClient session.BaseClient
}

func NewServerEventHandler(BaseServer session.BaseServer, BaseClient session.BaseClient) ServerEventHandler {
	return ServerEventHandler{
		BaseServer: BaseServer,
		BaseClient: BaseClient,
	}
}
