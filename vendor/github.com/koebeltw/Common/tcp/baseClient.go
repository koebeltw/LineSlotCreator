package tcp

import (
	"github.com/koebeltw/Common/type"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	stateNone byte = iota
	stateConnect
	stateDisConnect
	stateReConnect
)

type BaseClient interface {
	Connect() error
	Send(byte, byte, []byte)
	Close()
	IsConnect() bool
	SetUserHandler(UserHandler)
	SetCoder(Coder)
	SetReConnectSecond(time.Duration)
	SetEventHandler(*EventHandler)
	SetServerAddr(value Type.Addr)

	GetSession() (Session)
}

type tcpClientState struct {
	remoteAddr string
	remotePort int
	connected  bool
}

// Client TCP客户端描述
type baseClient struct {
	tcpClientState
	session         Session
	wg              sync.WaitGroup
	state           byte
	reConnectSecond time.Duration
	closeCh         chan bool
	//addr            string
	//port            int
	serverAddr   Type.Addr
	userHandler  UserHandler
	coder        Coder
	eventHandler *EventHandler
	//pool            *sync.Pool
}

// CreateBaseClient 创建一个TCPClient实例
func CreateBaseClient() BaseClient { return &baseClient{} }

func (c *baseClient) SetUserHandler(value UserHandler) {
	c.userHandler = value
}

func (c *baseClient) SetCoder(value Coder) {
	c.coder = value
}

func (c *baseClient) SetReConnectSecond(value time.Duration) {
	c.reConnectSecond = value
}

func (c *baseClient) SetEventHandler(value *EventHandler) {
	c.eventHandler = value
}

func (c *baseClient) SetServerAddr(value Type.Addr) {
	c.serverAddr = value
}

func (c *baseClient) GetSession() (Session) {
	return c.session
}

// Connect 连接到服务器
func (c *baseClient) Connect() error {
	//pool := &sync.Pool{New:func()interface{}{return packet.NewPacket()}}
	isReConn := false
	var reConnectSecond = time.Second * 3
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.serverAddr.IP, c.serverAddr.Port))
		if err != nil {
			fmt.Printf("net.Dial Error: %s\n", err)
			continue
		}

		c.session = nil
		c.session = NewSession(conn, c.userHandler, c.coder, &c.wg, isReConn, c.eventHandler)
		c.session.Start()
		//c.tcpClientState = tcpClientState{
		//	remoteAddr: c.addr,
		//	remotePort: c.port,
		//	connected:  true,
		//}

		// log.Printf("%s Connect Server %s ", C.session.LocalAddr(), C.session.RemoteAddr())

		c.wg.Wait()

		// log.Printf("%s Disconnect Server %s ", C.session.LocalAddr(), C.session.RemoteAddr())
		// log.Println("Session: " + strconv.FormatUint(uint64(C.session.ID), 10) + " disconnect")
		if c.reConnectSecond == 0 {
			break
		}
		isReConn = true
		reConnectSecond = c.reConnectSecond
		time.Sleep(reConnectSecond)
	}

	return nil
}

// Send 发送数据
func (c *baseClient) Send(msgNo byte, subNo byte, buffer []byte) {
	c.session.SendMsg(msgNo, subNo, buffer)
}

// Close 关闭连接
func (c *baseClient) Close() {
	c.session.Close()
}

// IsConnect 关闭连接
func (c *baseClient) IsConnect() bool {
	return c.session != nil
}
