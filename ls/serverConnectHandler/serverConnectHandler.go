package serverConnectHandler

import (
	"github.com/koebeltw/Common/session"
	"github.com/koebeltw/LineSlotCreator/ls/loadData"
	"log"
)

func NewServerConnectHandler(BaseServer session.BaseServer) *ServerConnectHandler {
	return &ServerConnectHandler{
		BaseServer: BaseServer,
	}
}

//// DBMgr blabla
type ServerConnectHandler struct{
	BaseServer session.BaseServer
	BaseClient session.BaseClient
}

// OnUserConnect 客户端连接事件
func (ser *ServerConnectHandler) OnUserConnect(s session.Session) {
	loadData.GetDataArrayBytes(func (b []byte){
		s.SendMsg(000, 001, b)
	})

	log.Printf("LS [%d] %s connect GS %s", s.GetID(), s.LocalAddr(), s.RemoteAddr())
}

// OnUserDisconnect 客户端断开连接事件
func (ser *ServerConnectHandler) OnUserDisconnect(s session.Session) {
	log.Printf("LS [%d] %s disconnect GS %s", s.GetID(), s.LocalAddr(), s.RemoteAddr())
}

// OnServerInit blabla
func (ser *ServerConnectHandler) OnServerInit(s session.BaseServer) {
	log.Println("OnServerInit")
}

// OnServerDestroy blabla
func (ser *ServerConnectHandler) OnServerDestroy(s session.BaseServer) {
	log.Println("OnServerDestroy")
}
