package gonet

import (
	"github.com/flylib/goutils/sync/spinlock"
	"sync/atomic"
)

type ISessionIdentify interface {
	ID() uint64
	SetID(id uint64)
	UpdateID(id uint64)
	WithContext(c *AppContext)
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
