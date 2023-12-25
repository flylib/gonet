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
	ctx *Context
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
	value, ok := s.ctx.sessions.alive.Load(s.id)
	if ok {
		s.ctx.sessions.alive.Delete(s.id)
		s.id = id
		s.ctx.sessions.alive.Store(s.id, value)
	}
}

func (s *SessionCommon) WithContext(c *Context) {
	s.ctx = c
}

func (s *SessionCommon) GetContext() *Context {
	return s.ctx
}

func (s *SessionCommon) Clear() {
	s.ctx = nil
	s.id = 0
	s.Store(&zeroData)
}
