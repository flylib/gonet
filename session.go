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
	//场景通道
	SessionSceneChan struct {
		messageQueues []chan Msg //消息队列
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
	newSession := sessions.Get()
	atomic.AddUint64(&sessions.autoIncrement, 1) //++
	newSession.(interface{ setID(id uint64) }).setID(sessions.autoIncrement)
	sessions.Store(sessions.autoIncrement, newSession)
	session := newSession.(Session)
	session.JoinOrUpdateActor(DefaultSceneID, defaultActor)
	HandleEvent(Event{
		Actor: defaultActor,
		context: context{
			session: session,
			data:    &msgSessionConnect,
		}})
	return session
}

func RecycleSession(s Session) {
	s.Close()
	sessions.Delete(s.ID())
	sessions.Put(s)
	HandleEvent(Event{
		Actor: defaultActor,
		context: context{
			session: s,
			data:    &msgSessionClose,
		}})
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

//增加场景消息订阅
func (s *SessionSceneChan) AddSceneSubscribe(sceneID uint8, ch <-chan Msg) {
	ch := make(chan Msg)
}

//@Param Actor id
func (s *SessionSceneChan) GetActor(ActorID int) (Actor, error) {
	if ActorID >= len(s.Actor) || s.Actor[ActorID] == nil {
		return nil, ErrNotFoundActor
	}
	return s.Actor[ActorID], nil
}

func RegisterSessionType(ses interface{}) {
	sessionType = reflect.TypeOf(ses)
}
