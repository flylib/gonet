package goNet

import (
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	sessionMgr  = NewSessionManager()
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
		// ID
		ID() uint64
		//数据存储
		Value(obj ...interface{}) interface{}
		//添加场景,如果场景相同会进行覆盖
		JoinScene(sceneID uint8, scene Scene)
		//获取场景
		GetScene(sceneID uint8) Scene
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
	//会话当前所在场景
	SessionScene struct {
		scenes []Scene
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

func GetSession(id uint64) (Session, bool) {
	value, ok := sessionMgr.Load(id)
	if ok {
		return value.(Session), ok
	}
	return nil, false
}

func AddSession() Session {
	newSession := sessionMgr.Get()
	atomic.AddUint64(&sessionMgr.autoIncrement, 1) //++
	newSession.(interface{ setID(id uint64) }).setID(sessionMgr.autoIncrement)
	sessionMgr.Store(sessionMgr.autoIncrement, newSession)
	session := newSession.(Session)
	//notify
	msg := &Msg{
		Session: session,
		SceneID: GetMsgSceneID(MsgIDSessionConnect),
		ID:      MsgIDSessionConnect,
		Data:    &msgSessionConnect,
	}
	PushWorkerPool(msg)
	return session
}

func RecycleSession(s Session) {
	s.Close()
	sessionMgr.Delete(s.ID())
	sessionMgr.Put(s)
	//notify
	msg := &Msg{
		Session: s,
		SceneID: GetMsgSceneID(MsgIDSessionConnect),
		ID:      MsgIDSessionConnect,
		Data:    &msgSessionClose,
	}
	PushWorkerPool(msg)
}

func SessionCount() int {
	sum := 0
	sessionMgr.Range(func(key, value interface{}) bool {
		sum++
		return true
	})
	return sum
}

//广播
func Broadcast(msg interface{}) {
	sessionMgr.Range(func(_, item interface{}) bool {
		item.(Session).Send(msg)
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
func (s *SessionScene) JoinScene(sceneID uint8, scene Scene) {
	if s.scenes == nil {
		s.scenes = make([]Scene, int(sceneID)+1)
	}
	more := sceneID + 1 - uint8(len(s.scenes))
	for i := uint8(0); i < more; i++ {
		s.scenes = append(s.scenes, nil)
	}
	s.scenes[sceneID] = scene
}

//增加场景消息订阅
func (s *SessionScene) GetScene(sceneID uint8) Scene {
	if uint8(len(s.scenes)) <= sceneID {
		return nil
	}
	return s.scenes[sceneID]
}

func SetSessionType(ses interface{}) {
	sessionType = reflect.TypeOf(ses)
}
