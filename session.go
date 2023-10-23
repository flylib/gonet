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
	//get the working Contentx
	Context() *AppContext
}

type ISessionIdentify interface {
	ID() uint64
	SetID(id uint64)
	UpdateID(id uint64)
	WithContext(c *AppContext)
	//get the working Contentx
	Context() *AppContext
	ClearIdentify()
}

// 核心会话标志
type SessionIdentify struct {
	*AppContext
	id uint64
}

func (s *SessionIdentify) ID() uint64 {
	return s.id
}

func (s *SessionIdentify) SetID(id uint64) {
	s.id = id
}

func (s *SessionIdentify) UpdateID(id uint64) {
	value, ok := s.AppContext.sessionMgr.alive.Load(s.id)
	if ok {
		s.AppContext.sessionMgr.alive.Delete(s.id)
		s.id = id
		s.AppContext.sessionMgr.alive.Store(s.id, value)
	}
}

func (s *SessionIdentify) WithContext(c *AppContext) {
	s.AppContext = c
}

func (s *SessionIdentify) Context() *AppContext {
	return s.AppContext
}

func (s *SessionIdentify) ClearIdentify() {
	s.AppContext = nil
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

func (s *SessionAbility) ClearAbility() {
	s.Store(nil)
}
