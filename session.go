package gonet

import "sync"

///////////////////////////////
/////    Session POOL   //////
//////////////////////////////

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

func (s SessionIdentify) ID() uint64 {
	return s.id
}
