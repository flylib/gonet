package gonet

import (
	"net"
	"sync"
)

///////////////////////////////
/////    Connection POOL   //////
//////////////////////////////

//会话
type Connection interface {
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
	ConnIdentify struct {
		//id
		id uint64
	}
	//存储功能
	ConnStore struct {
		sync.Map
	}
)

func (s *ConnIdentify) ID() uint64 {
	return s.id
}

func (s *ConnIdentify) setID(id uint64) {
	s.id = id
}
