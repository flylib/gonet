package goNet

import (
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	sessions    = NewSessionManager()
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
	}
)

func NewSessionManager() *sessionManager {
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
	session.JoinOrUpdateRoute(DefaultRouteID, defaultRoute)
	HandleEvent(Event{from: session, route: defaultRoute, data: &msgSessionConnect})
	return session
}

func RecycleSession(s Session) {
	s.Close()
	sessions.Delete(s.ID())
	sessions.Put(s)
	HandleEvent(Event{from: s, route: defaultRoute, data: &msgSessionClose})
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
		return nil, ErrNotFoundRoute
	}
	return s.route[routeID], nil
}

func RegisterSessionType(ses interface{}) {
	sessionType = reflect.TypeOf(ses)
}
