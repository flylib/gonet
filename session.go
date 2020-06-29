package goNet

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/panjf2000/ants/v2"
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

		//加入或者更新消息控制模块
		JoinOrUpdateController(index int, c Route)
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
		autoIncrement uint64 //流水号
		sync.Map
		*sync.Pool
		handler *ants.Pool //消息处理线程
	}
	//事件
	Event struct {
		session    Session     //会话
		controller Route       //控制器
		msg        interface{} //消息
	}
)

//构造事件
func CreateEvent(s Session, c Route, msg interface{}) Event {
	return Event{
		session:    s,
		controller: c,
		msg:        msg,
	}
}

func NewSessionManager() *sessionManager {
	handlers, _ := ants.NewPool(Default_Handler_Count)
	return &sessionManager{
		Pool: &sync.Pool{New: func() interface{} {
			return reflect.New(sessionType).Interface()
		}},
		handler: handlers,
	}
}

//处理事件
func HandleEvent(e Event) {
	if err := sessions.handler.Submit(func() {
		e.controller.OnMsg(e.session, e.msg)
	}); err != nil {
		logs.Error("antsPool commit message error,reason is ", err.Error())
	}
}

//调整处理者的数量
func TuneHandlerCount(count int) {
	sessions.handler.Tune(count)
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
	session.JoinOrUpdateController(System_Route_ID, sysRoute)
	//notify session connect
	HandleEvent(CreateEvent(session, sysRoute, &msgSessionConnect))
	return session
}

func RecycleSession(s Session) {
	s.Close()
	sessions.Delete(s.ID())
	sessions.Put(s)
	HandleEvent(CreateEvent(s, sysRoute, &msgSessionClose))
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

func (s *SessionRoute) JoinOrUpdateController(index int, c Route) {
	if index < 0 {
		return
	}
	if s.route == nil {
		s.route = make([]Route, 0, 3)
	}
	more := index - len(s.route) + 1
	//extend
	if more > 0 {
		moreControllers := make([]Route, more)
		s.route = append(s.route, moreControllers...)
	}
	s.route[index] = c
}

func (s *SessionRoute) GetController(index int) (Route, error) {
	if index >= len(s.route) || s.route[index] == nil {
		return nil, errors.New("not found route")
	}
	return s.route[index], nil
}

func RegisterSessionType(ses interface{}) {
	sessionType = reflect.TypeOf(ses)
}
