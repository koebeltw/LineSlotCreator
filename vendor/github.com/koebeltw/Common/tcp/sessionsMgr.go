package tcp

import (
	"github.com/koebeltw/Common/orderArrayManager"
	"github.com/gogf/gf/g/container/gchan"
	"github.com/koebeltw/Common/type"
)

type SessionMgr interface {
	Add(c Session) (index uint16, err error)
	Del(index uint16) (err error)
	Put(index uint16, c Session) (err error)
	Start()
	Clear()
	Close()
	SendMsgToAll(msgNo byte, subNo byte, buffer []byte)
	SendMsgToSomeOne(index uint16, msgNo byte, subNo byte, buffer []byte)
	SendMsgExclude(excludeIndex []uint16, msgNo byte, subNo byte, buffer []byte)
	FindSessionByID(id uint16) (p Session)
	IsNil(id uint16) bool
}

// SessionsMgr blabla
type sessionsMgr struct {
	sessions *orderArrayManager.OAM
	sendCh   *gchan.Chan
	//sendCh   chan func()
}

// NewSessionsMgr blabla
func NewSessionsMgr(max uint64) (result SessionMgr) {
	result = &sessionsMgr{sessions: orderArrayManager.NewOAM(max), sendCh: gchan.New(1000)}
	result.Start()
	return
}

// Add blabla
func (s *sessionsMgr) Add(c Session) (index uint16, err error) {
	i, e := s.sessions.Add(c)
	return uint16(i), e
}

// Del blabla
func (s *sessionsMgr) Del(index uint16) (err error) {
	return s.sessions.Del(uint64(index))
}

// Put blabla
func (s *sessionsMgr) Put(index uint16, c Session) (err error) {
	return s.sessions.Put(uint64(index), c)
}

// Clear blabla
func (s *sessionsMgr) Clear() {
	s.sessions.Clear()
}

// SendMsgToAll blabla
func (s *sessionsMgr) SendMsgToAll(msgNo byte, subNo byte, buffer []byte) {
	s.sendCh.Push(func() {
		s.sessions.ForEach(
			func(key interface{}, value interface{}) {
				value.(Session).SendMsg(msgNo, subNo, buffer)
			})
	})
}

// SendMsgToSomeOne blabla
func (s *sessionsMgr) SendMsgToSomeOne(index uint16, msgNo byte, subNo byte, buffer []byte) {
	s.sendCh.Push(func() {
		s.sessions.Find(
			func(key interface{}, value interface{}) bool {
				session := value.(Session)
				if session.GetID() == uint16(index) {
					session.SendMsg(msgNo, subNo, buffer)
					return true
				}

				return false
			})
	})
}

// SendMsgExclude blabla
func (s *sessionsMgr) SendMsgExclude(excludeIndex []uint16, msgNo byte, subNo byte, buffer []byte) {
	s.sendCh.Push(func() {
		s.sessions.ForEach(
			func(key interface{}, value interface{}) {
				session := value.(Session)
				for index := range excludeIndex {
					if session.GetID() == uint16(index) {
						return
					}
				}

				session.SendMsg(msgNo, subNo, buffer)
			})
	})
}

// FindSessionByID blabla
func (s *sessionsMgr) FindSessionByID(id uint16) (p Session) {
	return s.sessions.Find(
		func(key interface{}, value interface{}) bool {
			if session, ok := value.(Session); ok {
				if session.GetID() == uint16(id) {
					return true
				}
			}

			return false
		}).(Session)
}

// FindSessionByID blabla
func (s *sessionsMgr) IsNil(id uint16) bool {
	return s.sessions.IsNil(uint64(id))
}

func (s *sessionsMgr) Close() {
	s.sessions.ForEach(
		func(key interface{}, value interface{}) {
			session := value.(Session)
			session.Close()
		})

	s.sessions.Clear()
}

func (s *sessionsMgr) Start()  {
	go s.sendThread()
}

func (s *sessionsMgr) sendThread() {
	for {
		if v := s.sendCh.Pop(); v != nil {
			if v, ok := v.(Type.Function); ok {
				v()
			}
		}
	}
}
