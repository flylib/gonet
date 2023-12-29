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
}

// 核心会话标志
type SessionCommon struct {
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
	value, ok := defaultCtx.sessionManager.alive.Load(s.id)
	if ok {
		defaultCtx.sessionManager.alive.Delete(s.id)
		s.id = id
		defaultCtx.sessionManager.alive.Store(s.id, value)
	}
}
func (s *SessionCommon) Clear() {
	s.id = 0
	s.Value = atomic.Value{}
}
