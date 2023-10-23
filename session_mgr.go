package gonet

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// 会话管理
type SessionManager struct {
	aliveNum uint32    //当前活跃链接总数
	incr     uint64    //流水号
	alive    sync.Map  //活跃链接
	idle     sync.Pool //空闲会话
}

func NewSessionManager(sessionType reflect.Type) *SessionManager {
	return &SessionManager{
		idle: sync.Pool{New: func() any {
			return reflect.New(sessionType).Interface()
		}},
	}
}

// 活跃会话
func (s *SessionManager) AddAliveSession(session ISession) {
	atomic.AddUint32(&s.aliveNum, 1)
	session.(interface{ SetID(id uint64) }).SetID(atomic.AddUint64(&s.incr, 1))
	s.alive.Store(session.ID(), session)
}

func (s *SessionManager) GetAliveSession(id uint64) (session ISession, exist bool) {
	ss, ok := s.alive.Load(id)
	if !ok {
		return nil, ok
	}
	return ss.(ISession), ok
}

func (s *SessionManager) CountAliveSession() uint32 {
	return atomic.LoadUint32(&s.aliveNum)
}

func (s *SessionManager) GetIdleSession() ISession {
	return s.idle.Get().(ISession)
}
func (s *SessionManager) RecycleIdleSession(session ISession) {
	atomic.AddUint32(&s.aliveNum, -1)
	s.alive.Delete(session.ID())
	s.idle.Put(session)
}
