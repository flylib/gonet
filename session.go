package goNet

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	sessions    = NewSessionManager()
	sessionType reflect.Type
)

const (
	Default_Handler_Count = 10
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

		//加入或者更新路由
		JoinOrUpdateRoute(index int, c Route)
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
	SessionRoute struct {
		route []Route
	}
	//会话管理
	sessionManager struct {
		sync.Map
		*sync.Pool
		autoIncrement uint64 //流水号
		//handler       *ants.Pool //消息处理线程
	}
)

//构造事件
//func NewEvent(s Session, c Route, msg interface{}) Event {
//	return Event{
//		From:    s,
//		Handler: c,
//		Msg:     msg,
//	}
//}

func NewSessionManager() *sessionManager {
	//handlers, _ := ants.NewPool(Default_Handler_Count)
	return &sessionManager{
		Pool: &sync.Pool{New: func() interface{} {
			return reflect.New(sessionType).Interface()
		}},
		//handler: handlers,
		//ch:      make(chan Event),
	}
}

//处理事件
//func HandleEvent(e Event) {
//	if err := sessions.handler.Submit(func() {
//		e.controller.OnMsg(e.session, e.msg)
//	}); err != nil {
//		logs.Error("antsPool commit message error,reason is ", err.Error())
//	}
//}

//调整处理者的数量
//func TuneHandlerCount(count int) {
//	sessions.handler.Tune(count)
//}

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
	session.JoinOrUpdateRoute(System_Route_ID, sysRoute)
	//notify session connect
	//HandleEvent(NewEvent(session, sysRoute, &msgSessionConnect))
	CommitWorkerPool(Event{From: session, Router: sysRoute, Msg: &msgSessionConnect})
	return session
}

func RecycleSession(s Session) {
	s.Close()
	sessions.Delete(s.ID())
	sessions.Put(s)
	//HandleEvent(NewEvent(s, sysRoute, &msgSessionClose))
	CommitWorkerPool(Event{From: s, Router: sysRoute, Msg: &msgSessionClose})
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

func (s *SessionRoute) JoinOrUpdateRoute(id int, c Route) {
	if id < 0 {
		return
	}
	if s.route == nil {
		s.route = make([]Route, 0, 3)
	}
	more := id - len(s.route) + 1
	//extend
	if more > 0 {
		moreControllers := make([]Route, more)
		s.route = append(s.route, moreControllers...)
	}
	s.route[id] = c
}

//@Param route id
func (s *SessionRoute) GetRoute(routeID int) (Route, error) {
	if routeID >= len(s.route) || s.route[routeID] == nil {
		return nil, errors.New("not found route")
	}
	return s.route[routeID], nil
}

func RegisterSessionType(ses interface{}) {
	sessionType = reflect.TypeOf(ses)
}
