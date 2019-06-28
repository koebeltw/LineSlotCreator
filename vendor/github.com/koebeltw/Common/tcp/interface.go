package tcp

// userHandler blabla
type UserHandler interface {
	OnUserConnect(s Session)
	OnUserDisconnect(s Session)
	//OnUserReConnect(s Session)
}

// serverHandler blabla
type ServerHandler interface {
	OnServerInit(s BaseServer)
	OnServerDestroy(s BaseServer)
}

// coder blabla
type Coder interface {
	Decode(c Session) (r EventMsg, err error)
	Encode(c Session, msgNo byte, subNo byte, buffer []byte) (r []byte, err error)
}
