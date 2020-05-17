package goNet

import (
	"errors"
	"reflect"
	"sync"
)

var (
	SessionManager = newSessionManager()
	sessionType    reflect.Type
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

		// ID b
		ID() uint64

		//数据存储
		Value(obj ...interface{}) interface{}
	}
	//核心会话标志
	SessionIdentify struct {
		//id
		id uint64
	}
	//存储功能
	SessionStore struct {
		obj interface{}
	}
	//消息路由
	SessionController struct {
		//example center_service/room_service/...
		controllers []Controller
	}

	//session管理器
	sessionManager struct {
		//流水号
		AutoIncrement uint64
		//活跃sessions
		sessions map[uint64]Session
		*sync.Pool
		//保证线程安全
		sync.RWMutex
	}
)

func newSessionManager() *sessionManager {
	return &sessionManager{
		//idleSessions: map[uint32]Session{},
		sessions: map[uint64]Session{},
		Pool: &sync.Pool{New: func() interface{} {
			return reflect.New(sessionType).Interface()
		}},
	}
}

func (s *sessionManager) GetSessionById(id uint64) Session {
	return s.sessions[id]
}

func (s *sessionManager) AddSession(ses Session) {
	s.sessions[ses.ID()] = ses
	ses.(interface{ JoinController(index int, c Controller) }).JoinController(SYSTEM_CONTROLLER_IDX, systemController)
	//notify session connect
	SubmitMsgToAntsPool(systemController, ses, &SessionConnect{})
}

//回收到空闲会话池
func (s *sessionManager) RecycleSession(ses Session) {
	ses.Close()
	delete(s.sessions, ses.ID())
	s.Put(ses)
	//notify session close
	SubmitMsgToAntsPool(systemController, ses, &SessionClose{})
}

//总数
func (s *sessionManager) GetSessionCount() int {
	return len(s.sessions)
}

//广播
func (s *sessionManager) Broadcast(msg interface{}) {
	for _, ses := range s.sessions {
		ses.Send(msg)
	}
}

func (s *SessionIdentify) ID() uint64 {
	return s.id
}

func (s *SessionIdentify) SetID(id uint64) {
	s.id = id
}

func (s *SessionStore) Value(v ...interface{}) interface{} {
	if len(v) > 0 {
		s.obj = v[0]
	}
	return s.obj
}

func (s *SessionController) JoinController(index int, c Controller) {
	if s.controllers == nil {
		s.controllers = make([]Controller, 0, 3)
	}
	more := index - len(s.controllers) + 1
	//extend
	if more > 0 {
		moreControllers := make([]Controller, more)
		s.controllers = append(s.controllers, moreControllers...)
	}
	s.controllers[index] = c
}

func (s *SessionController) GetController(index int) (Controller, error) {
	if index >= len(s.controllers) || s.controllers[index] == nil {
		return nil, errors.New("not found controller")
	}
	return s.controllers[index], nil
}

func RegisterSessionType(ses interface{}) {
	sessionType = reflect.TypeOf(ses)
}
