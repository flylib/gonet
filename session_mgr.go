package gonet

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// 会话管理
type SessionManager struct {
	incr  uint64    //流水号
	alive sync.Map  //活跃链接
	idle  sync.Pool //空闲会话
}

func newSessionManager(sessionType reflect.Type) *SessionManager {
	return &SessionManager{
		idle: sync.Pool{New: func() any {
			return reflect.New(sessionType).Interface()
		}},
	}
}

// 活跃会话
func (self *SessionManager) addAliveSession(session ISession) {
	session.(interface{ SetID(id uint64) }).SetID(atomic.AddUint64(&self.incr, 1))
	self.alive.Store(session.ID(), session)
}

func (self *SessionManager) getAliveSession(id uint64) (session ISession, exist bool) {
	s, ok := self.alive.Load(id)
	if !ok {
		return nil, ok
	}
	return s.(ISession), ok
}

func (self *SessionManager) removeAliveSession(session ISession) {
	session.Close()
	session.(interface{ Clear() }).Clear()
	self.alive.Delete(session.ID())
	self.recycleIdleSession(session)
}

func (self *SessionManager) CountAliveSession() int {
	total := 0
	self.alive.Range(func(key, value any) bool {
		total++
		return true
	})
	return total
}

// 空闲会话
func (self *SessionManager) getIdleSession() ISession {
	return self.idle.Get().(ISession)
}
func (self *SessionManager) recycleIdleSession(session ISession) {
	self.alive.Delete(session.ID())
	self.idle.Put(session)
}
