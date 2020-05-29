package goNet

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	sessions    = newSessionManager()
	sessionType reflect.Type
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

		//加入或者更新消息控制模块
		JoinOrUpdateController(index int, c Controller)
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
	//会话管理
	sessionManager struct {
		autoIncrement uint64 //流水号
		sync.Map
		*sync.Pool
	}
)

func newSessionManager() *sessionManager {
	return &sessionManager{
		Pool: &sync.Pool{New: func() interface{} {
			return reflect.New(sessionType).Interface()
		}},
	}
}

func FindSession(id uint64) (Session, bool) {
	value, ok := sessions.Load(id)
	if ok {
		return value.(Session), ok
	}
	return nil, false
}

func AddSession() Session {
	ses := sessions.Get()
	atomic.AddUint64(&sessions.autoIncrement, 1)
	ses.(interface{ setID(id uint64) }).setID(sessions.autoIncrement)
	sessions.Store(sessions.autoIncrement, ses)
	session := ses.(Session)
	session.JoinOrUpdateController(SYSTEM_CONTROLLER_IDX, systemController)
	//notify session connect
	SubmitMsgToAntsPool(systemController, session, &msgSessionConnect)
	return session
}

func RecycleSession(ses Session) {
	ses.Close()
	sessions.Delete(ses.ID())
	sessions.Put(ses)
	SubmitMsgToAntsPool(systemController, ses, &msgSessionClose)
}

func SessionCount() int {
	sum := 0
	sessions.Range(func(key, value interface{}) bool {
		sum++
		return true
	})
	return sum
}

//广播
func Broadcast(msg interface{}) {
	sessions.Range(func(_, value interface{}) bool {
		value.(Session).Send(msg)
		return true
	})
}

func (s *SessionIdentify) ID() uint64 {
	return s.id
}

func (s *SessionIdentify) setID(id uint64) {
	s.id = id
}

func (s *SessionStore) Value(v ...interface{}) interface{} {
	if len(v) > 0 {
		s.obj = v[0]
	}
	return s.obj
}

func (s *SessionController) JoinOrUpdateController(index int, c Controller) {
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
