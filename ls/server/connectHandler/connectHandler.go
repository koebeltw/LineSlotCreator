package connectHandler

import (
	"github.com/koebeltw/Common/tcp"
	"github.com/koebeltw/LineSlotCreator/ls/loadData"
	"log"
)

func NewConnectHandler(BaseServer tcp.BaseServer) *ConnectHandler {
	return &ConnectHandler{
		BaseServer: BaseServer,
	}
}

type ConnectHandler struct{
	BaseServer tcp.BaseServer
	BaseClient tcp.BaseClient
}

// OnUserConnect 客户端连接事件
func (ser *ConnectHandler) OnUserConnect(s tcp.Session) {
	loadData.GetDataArrayBytes(func (b []byte){
		s.SendMsg(000, 001, b)
	})

	log.Printf("[%d] %s connect %s", s.GetID(), s.LocalAddr(), s.RemoteAddr())
}

// OnUserDisconnect 客户端断开连接事件
func (ser *ConnectHandler) OnUserDisconnect(s tcp.Session) {
	log.Printf("[%d] %s disconnect %s", s.GetID(), s.LocalAddr(), s.RemoteAddr())
}

// OnServerInit blabla
func (ser *ConnectHandler) OnServerInit(s tcp.BaseServer) {
	log.Println("OnServerInit")
}

// OnServerDestroy blabla
func (ser *ConnectHandler) OnServerDestroy(s tcp.BaseServer) {
	log.Println("OnServerDestroy")
}
