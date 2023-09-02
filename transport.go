package gonet

import (
	"net"
	"reflect"
	"sync"
)

type TransportProtocol string

const (
	TCP  TransportProtocol = "tcp"
	KCP  TransportProtocol = "kcp"
	UDP  TransportProtocol = "udp"
	WS   TransportProtocol = "websocket"
	HTTP TransportProtocol = "http"
	QUIC TransportProtocol = "quic"
	RPC  TransportProtocol = "rpc"
)

// Interfaces
type (
	//服务端
	IServer interface {
		// 启动监听
		Listen() error
		// 停止服务
		Stop() error
		// 地址
		Addr() string
		// 会话类型
		SessionType() reflect.Type
	}
	//客户端
	IClient interface {
	}
	//会话
	ISession interface {
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
)

// server端属性
type ServerIdentify struct {
	uuid string
	//地址
	addr string
}

func (s *ServerIdentify) Addr() string {
	return s.addr
}

func (s *ServerIdentify) SetAddr(addr string) {
	s.addr = addr
}

// 存储功能
type SessionStore struct {
	sync.Map
}

// 核心会话标志
type SessionIdentify struct {
	id uint64
}

func (s *SessionIdentify) ID() uint64 {
	return s.id
}

func (s *SessionIdentify) SetID(id uint64) {
	s.id = id
}
