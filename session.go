package goNet

import (
	"sync"
	"sync/atomic"
)

var (
	SessionManager = newSessionManger()
)

type (
	//会话
	Session interface {
		//原始套接字
		Socket() interface{}

		// 发送消息，消息需要以指针格式传入
		Send(msg interface{})

		// 断开
		Close()

		// ID
		ID() uint32

		//数据存储
		Value(v ...interface{}) interface{}
	}
	//核心会话标志
	SessionIdentify struct {
		//id
		id uint32
	}
	//session管理器
	sessionManager struct {
		//流水号
		count uint32
		//空闲会话，重复利用
		idleSessions map[uint32]Session
		//活跃sessions
		sessions map[uint32]Session
		//保证线程安全
		sync.RWMutex
		//会话断开消息
		OnSessionClose func(session Session)
	}
)

func newSessionManger() *sessionManager {
	return &sessionManager{
		idleSessions: map[uint32]Session{},
		sessions:     map[uint32]Session{},
	}
}

func (s *sessionManager) GetIdleSession() Session {
	s.Lock()
	defer s.Unlock()
	for _, ses := range s.idleSessions {
		delete(s.idleSessions, ses.ID())
		return ses
	}
	return nil
}

func (s *sessionManager) GetSessionById(id uint32) Session {
	s.RLock()
	defer s.RUnlock()

	return s.sessions[id]
}

func (s *sessionManager) AddSession(ses Session) {
	s.Lock()
	defer s.Unlock()

	//如果会话新分配的要分配流水号
	if ses.ID() < 1 {
		s.count++
		ses.(interface {
			SetID(uint32)
		}).SetID(s.count)
	}
	s.sessions[ses.ID()] = ses
}

//回收到空闲会话池
func (s *sessionManager) RecycleSession(ses Session) {
	s.Lock()
	defer s.Unlock()

	if s.OnSessionClose != nil {
		s.OnSessionClose(ses)
	}
	ses.Close()
	delete(s.sessions, ses.ID())
	s.idleSessions[ses.ID()] = ses
}

//总数
func (s *sessionManager) GetSessionCount() uint32 {
	return atomic.LoadUint32(&s.count)
}

//活跃总数
func (s *sessionManager) GetSessionAliveCount() uint32 {
	s.RLock()
	defer s.RUnlock()
	return uint32(len(s.sessions))
}

//广播
func (s *sessionManager) Broadcast(msg interface{}) {
	for _, ses := range s.sessions {
		ses.Send(msg)
	}
}

func (s *SessionIdentify) ID() uint32 {
	return s.id
}

func (s *SessionIdentify) SetID(id uint32) {
	s.id = id
}
