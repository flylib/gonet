package gonet

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// 会话管理
type SessionManager struct {
	aliveNum int32     //当前活跃链接总数
	incr     uint64    //流水号
	alive    sync.Map  //活跃链接
	idle     sync.Pool //空闲会话
}

func newSessionManager(sessionType reflect.Type) *SessionManager {
	return &SessionManager{
		idle: sync.Pool{New: func() any {
			return reflect.New(sessionType).Interface()
		}},
	}
}

// 活跃会话
func (self *SessionManager) AddAliveSession(session ISession) {
	atomic.AddInt32(&self.aliveNum, 1)
	session.(interface{ SetID(id uint64) }).SetID(atomic.AddUint64(&self.incr, 1))
	self.alive.Store(session.ID(), session)
}

func (self *SessionManager) GetAliveSession(id uint64) (session ISession, exist bool) {
	s, ok := self.alive.Load(id)
	if !ok {
		return nil, ok
	}
	return s.(ISession), ok
}

func (self *SessionManager) CountAliveSession() int32 {
	return atomic.LoadInt32(&self.aliveNum)
}

// 空闲会话
func (self *SessionManager) GetIdleSession() ISession {
	return self.idle.Get().(ISession)
}
func (self *SessionManager) RecycleIdleSession(session ISession) {
	atomic.AddInt32(&self.aliveNum, -1)
	self.alive.Delete(session.ID())
	self.idle.Put(session)
}
