package tcp

import (
	"github.com/koebeltw/Common/type"
	"fmt"
	"log"
	"net"
	"sync"
	"strings"
)

type BaseServer interface {
	Start() error
	Stop()
	//SessionCount() uint32
	//SetMaxSession(int)
	Close()
	GetSessionMgr() SessionMgr
	SetSessionMgr(SessionMgr)
	SetUserHandler(UserHandler)
	SetServerHandler(ServerHandler)
	SetCoder(Coder)
	SetEventHandler(*EventHandler)
	SetServerAddr(value Type.Addr)

	GetClientsAddr() []Type.Addr
	SetClientsAddr(value []Type.Addr)

	OnUserConnect(s Session)
	OnUserDisconnect(s Session)
	SendMsgToAll(msgNo byte, subNo byte, buffer []byte)
	SendMsgToSomeOne(index uint16, msgNo byte, subNo byte, buffer []byte)
	SendMsgExclude(excludeIndex []uint16, msgNo byte, subNo byte, buffer []byte)
	FindSessionByID(id uint16) (p Session)
}

// Server 描述一个服务器的结构
type baseServer struct {
	// Sessions map[string]*Session
	//ServerState
	//MaxSessionCount int            // 最大连接数，为0则不限制服务器最大连接数
	listener net.Listener // 监听句柄
	//terminated bool           // 通知是否停止Service
	wg sync.WaitGroup // 等待所有goroutine结束
	//addr          string
	//port          int
	//Addr          Type.Addr
	SessionsMgr   SessionMgr
	serverHandler ServerHandler
	userHandler   UserHandler
	coder         Coder
	eventHandler  *EventHandler
	data          interface{}
	serverAddr    Type.Addr
	clientsAddr    []Type.Addr
}

func (s *baseServer) SetUserHandler(value UserHandler) {
	s.userHandler = value
}

func (s *baseServer) SetServerHandler(value ServerHandler) {
	s.serverHandler = value
}

func (s *baseServer) SetCoder(value Coder) {
	s.coder = value
}

func (s *baseServer) SetEventHandler(value *EventHandler) {
	s.eventHandler = value
}

func (s *baseServer) GetSessionMgr() SessionMgr {
	return s.SessionsMgr
}

func (c *baseServer) SetServerAddr(value Type.Addr) {
	c.serverAddr = value
}

func (c *baseServer) GetClientsAddr() []Type.Addr {
	return c.clientsAddr
}

func (c *baseServer) SetClientsAddr(value []Type.Addr){
	c.clientsAddr = value
}

func (s *baseServer) SetSessionMgr(value SessionMgr) {
	s.SessionsMgr = value
}

func (s *baseServer) OnUserConnect(se Session){
	if s.clientsAddr != nil{
		isClose := true
		for _, value := range s.clientsAddr {
			connIP := strings.Split(se.GetConn().RemoteAddr().String(), ":")[0]
			if connIP == value.IP {
				if s.GetSessionMgr().IsNil(uint16(value.ID)){
					s.GetSessionMgr().Put(uint16(value.ID), se)
					se.SetID(uint16(value.ID))
					//session := NewSession(uint16(value.ID), se, s.userHandler, s.coder, &s.wg, false, s.eventHandler)
					isClose = false
				}
				break
			}
		}
		if isClose {
			se.Close()
			return
		}
	}else{
		//session := NewSession(conn, s.userHandler, s.coder, &s.wg, false, s.eventHandler)
		if id, err := s.GetSessionMgr().Add(se); err != nil{
			se.Close()
			return
		} else {
			s.GetSessionMgr().Put(uint16(id), se)
			se.SetID(uint16(id))
		}
	}

	if s.userHandler != nil {
		s.userHandler.OnUserConnect(se)
	}
}

func (s *baseServer) OnUserDisconnect(se Session){
	if !s.GetSessionMgr().IsNil(uint16(se.GetID())) {
		s.GetSessionMgr().Del(uint16(se.GetID()))
	}

	if s.userHandler != nil {
		s.userHandler.OnUserDisconnect(se)
	}
}

// CreateBaseServer 创建一个Server, 返回*Server
func CreateBaseServer() BaseServer { return &baseServer{} }

// Start 开始服务
func (s *baseServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.serverAddr.IP, s.serverAddr.Port))
	if err != nil {
		fmt.Printf("net.Listen Error: %s\n", err)
		return err
	}

	if s.serverHandler != nil {
		s.serverHandler.OnServerInit(s)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			//	if s.terminated {
			//		s.wg.Done()
			//		break
			//	}

			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				// S.processError(nil, err)
				break
			}

			NewSession(conn, s, s.coder, &s.wg, false, s.eventHandler).Start()
		}
	}()

	s.wg.Wait()
	if s.serverHandler != nil {
		s.serverHandler.OnServerDestroy(s)
	}
	log.Println("Server End")
	listener.Close()
	return nil
}

// Stop 停止服务
func (s *baseServer) Stop() {
	s.Close()
	s.listener.Close()
	//s.terminated = true
	// S.wg.Wait() // 等待结束
}

// Close blabla
func (s *baseServer) Close() {
	s.GetSessionMgr().Close()
}

func (s *baseServer) SendMsgToAll(msgNo byte, subNo byte, buffer []byte){
	s.SessionsMgr.SendMsgToAll(msgNo , subNo , buffer)
}

func (s *baseServer) SendMsgToSomeOne(index uint16, msgNo byte, subNo byte, buffer []byte){
	s.SessionsMgr.SendMsgToSomeOne(index, msgNo , subNo , buffer)
}

func (s *baseServer) SendMsgExclude(excludeIndex []uint16, msgNo byte, subNo byte, buffer []byte){
	s.SessionsMgr.SendMsgExclude(excludeIndex, msgNo , subNo , buffer)
}

func (s *baseServer) FindSessionByID(id uint16) (p Session){
	return s.SessionsMgr.FindSessionByID(id)
}
