package gonet

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type TransportType reflect.Type

//单次会话
type Session struct {
	Conn Conn //所属链接
	body interface{}
}

type (

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
)

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
	peerType = reflect.TypeOf(ses)
}
