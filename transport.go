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
		//会话类型
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
		Send(msg any) error
		//设置键值对，存储关联数据
		Store(key, value any)
		//获取键值对
		Load(key any) (value any, ok bool)
		//地址
		RemoteAddr() net.Addr
	}
)

// server端属性
type ServerIdentify struct {
	*Context
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
func (s *ServerIdentify) setContext(c *Context) {
	s.Context = c
}

// 存储功能
type SessionStore struct {
	sync.Map
}

func (s *SessionStore) Clear() {
	s.Range(func(key, value any) bool {
		s.Delete(key)
		return true
	})
}

// 核心会话标志
type SessionIdentify struct {
	*Context
	id uint64
}

func (self *SessionIdentify) ID() uint64 {
	return self.id
}

func (self *SessionIdentify) SetID(id uint64) {
	self.id = id
}

// 核心功能
type SessionAbility struct {
	once     *sync.Once
	wChannel chan any
}

func (s *SessionAbility) Init(size int) {
	if size < 1 {
		size = 1
	}
	s.wChannel = make(chan any, size)
	s.once = &sync.Once{}
}

func (s *SessionAbility) Close() {
	close(s.wChannel)
	s.wChannel = nil
	s.once = nil
}

func (s *SessionAbility) SendQueue(msg any) {
	s.wChannel <- msg
}

func (s *SessionAbility) WriteLoop(session ISession) {
	s.once.Do(func() {
		go func() {
			for msg := range s.wChannel {
				session.Send(msg)
			}
		}()
	})
}
