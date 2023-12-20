package gonet

import (
	"github.com/flylib/goutils/sync/spinlock"
	"net"
	"sync/atomic"
)

type ISession interface {
	//ID
	ID() uint64
	//close the connection
	Close() error
	//send the message to the other side
	Send(msgID uint32, msg any) error
	//remote addr
	RemoteAddr() net.Addr
	//convenient session storage data
	Store(value any)
	//load the data
	Load() (value any)
	//get the working Context
	GetContext() *Context
}

// 核心会话标志
type SessionCommon struct {
	*Context
	atomic.Value
	spinlock.Locker
	id uint64
}

func (s *SessionCommon) ID() uint64 {
	return s.id
}

func (s *SessionCommon) SetID(id uint64) {
	s.id = id
}

func (s *SessionCommon) UpdateID(id uint64) {
	value, ok := s.Context.sessionMgr.alive.Load(s.id)
	if ok {
		s.Context.sessionMgr.alive.Delete(s.id)
		s.id = id
		s.Context.sessionMgr.alive.Store(s.id, value)
	}
}

func (s *SessionCommon) WithContext(c *Context) {
	s.Context = c
}

func (s *SessionCommon) GetContext() *Context {
	return s.Context
}

func (s *SessionCommon) Clear() {
	s.Context = nil
	s.id = 0
	s.Store(&zeroData)
}
