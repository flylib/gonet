package gonet

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// 会话管理
type SessionManager struct {
	aliveNum int32     //当前活跃链接总数
	serial   uint64    //流水号
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

func (s *SessionManager) AddSession(session ISession) {
	atomic.AddInt32(&s.aliveNum, 1)
	session.(interface{ SetID(id uint64) }).SetID(atomic.AddUint64(&s.serial, 1))
	s.alive.Store(session.ID(), session)
}

func (s *SessionManager) GetSession(id uint64) (session ISession, exist bool) {
	ss, ok := s.alive.Load(id)
	if !ok {
		return nil, ok
	}
	return ss.(ISession), ok
}

func (s *SessionManager) SumAliveSession() int32 {
	return atomic.LoadInt32(&s.aliveNum)
}

func (s *SessionManager) GetIdleSession() ISession {
	return s.idle.Get().(ISession)
}

func (s *SessionManager) RecycleSession(session ISession) {
	session.Close()
	session.(interface{ Clear() }).Clear()
	atomic.AddInt32(&s.aliveNum, -1)
	s.alive.Delete(session.ID())
	s.idle.Put(session)
}

func (s *SessionManager) Broadcast(msgId uint32, msg any) {
	s.alive.Range(func(_, item interface{}) bool {
		session, ok := item.(ISession)
		if ok {
			session.Send(msgId, msg)
		}
		return true
	})
}
