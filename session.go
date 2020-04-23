package goNet

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	SessionManager = newSessionManager()
)

const (
	//无效会话id
	INVALID_SESSION_ID uint32 = 0
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

	MsgControllers struct {
		//example center_service/room_service/...
		controllers []MsgController
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
	}
)

func newSessionManager() *sessionManager {
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
	if ses.ID() == INVALID_SESSION_ID {
		s.count++
		ses.(interface {
			SetID(uint32)
		}).SetID(s.count)
		ses.(interface {
			AddController(index int, c MsgController)
		}).AddController(SYSTEM_CONTROLLER_IDX, systemMsgController)
	}
	s.sessions[ses.ID()] = ses
	//notify session connect
	SubmitMsgToAntsPool(systemMsgController, ses, &SessionConnect{})
}

//回收到空闲会话池
func (s *sessionManager) RecycleSession(ses Session) {
	s.Lock()
	defer s.Unlock()
	ses.Close()
	delete(s.sessions, ses.ID())
	s.idleSessions[ses.ID()] = ses
	//notify session close
	SubmitMsgToAntsPool(systemMsgController, ses, &SessionClose{})
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

func (s *MsgControllers) AddController(index int, c MsgController) {
	if s.controllers == nil {
		s.controllers = make([]MsgController, 0, 3)
	}
	more := index - len(s.controllers) + 1
	//extend
	if more > 0 {
		moreControllers := make([]MsgController, more)
		s.controllers = append(s.controllers, moreControllers...)
	}
	s.controllers[index] = c
}

func (s *MsgControllers) GetController(index int) (MsgController, error) {
	if index >= len(s.controllers) || s.controllers[index] == nil {
		return nil, errors.New("not found controller")
	}
	return s.controllers[index], nil
}
