package gonet

import (
	"net"
	"sync"
)

/*----------------------------------------------------------------
			///////////////////////////////
			/////    Session POOL   //////
			//////////////////////////////
----------------------------------------------------------------*/

//会话
type Session interface {
	//ID
	ID() uint64
	//断开
	Close() error
	//发送消息
	Send(msg interface{}) error
	//设置键值对，存储关联数据
	Store(key, value interface{})
	//获取键值对
	Load(key interface{}) (value interface{}, ok bool)
	//地址
	RemoteAddr() net.Addr
}

type (
	//核心会话标志
	SessionIdentify struct {
		//id
		id uint64
	}
	//存储功能
	SessionStore struct {
		sync.Map
	}
)

func (s *SessionIdentify) ID() uint64 {
	return s.id
}

func (s *SessionIdentify) setID(id uint64) {
	s.id = id
}

//会话管理
type SessionManager struct {
	sync.RWMutex
	incr     uint64    //流水号
	sessions sync.Map  //所有链接
	pool     sync.Pool //临时对象池
}

func (s *SessionManager) store(id uint64, session interface{}) {
	session.(interface{ setID(id uint64) }).setID(id)
	s.sessions.Store(id, session)
}

func (s *SessionManager) del(id uint64) {
	s.sessions.Delete(id)
}
