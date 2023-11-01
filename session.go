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

type ISessionIdentify interface {
	ID() uint64
	SetID(id uint64)
	UpdateID(id uint64)
	WithContext(c *Context)
	//get the working Contentx
	Context() *Context
	ClearIdentify()
}

// 核心会话标志
type SessionIdentify struct {
	*Context
	id uint64
}

func (s *SessionIdentify) ID() uint64 {
	return s.id
}

func (s *SessionIdentify) SetID(id uint64) {
	s.id = id
}

func (s *SessionIdentify) UpdateID(id uint64) {
	value, ok := s.Context.sessionMgr.alive.Load(s.id)
	if ok {
		s.Context.sessionMgr.alive.Delete(s.id)
		s.id = id
		s.Context.sessionMgr.alive.Store(s.id, value)
	}
}

func (s *SessionIdentify) WithContext(c *Context) {
	s.Context = c
}

func (s *SessionIdentify) GetContext() *Context {
	return s.Context
}

func (s *SessionIdentify) ClearIdentify() {
	s.Context = nil
	s.id = 0
}

type ISessionAbility interface {
	ClearAbility()
}

// common ability
type SessionAbility struct {
	atomic.Value
	spinlock.Locker
}

type invalidData struct {
}

var zeroData = invalidData{}

func (s *SessionAbility) ClearAbility() {
	s.Store(&zeroData)
}
