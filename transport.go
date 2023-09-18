package gonet

import (
	"net"
	"reflect"
	"sync/atomic"
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
		Dial()
	}
	//会话
	ISession interface {
		//ID
		ID() uint64
		//断开
		Close() error
		//发送消息
		Send(msg any) error
		//地址
		RemoteAddr() net.Addr
		//设置键值对，存储关联数据
		Store(value any)
		//获取键值对
		Load() (value any, ok bool)
	}
	ISessionIdentify interface {
		ID() uint64
		ClearIdentify()
		SetID(id uint64)
		UpdateID(id uint64)
		WithContext(c *Context)
		IsClosed() bool
		SetClosedStatus()
	}
	ISessionAbility interface {
		Store(val any)
		Load() (val any, ok bool)
		InitSendChanel()
		WriteSendChannel(buf []byte)
		RunningSendLoop(handler func([]byte))
		StopAbility()
	}
)

var (
	_ ISessionIdentify = new(SessionIdentify)
	_ ISessionAbility  = new(SessionAbility)
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

func (s *ServerIdentify) WithContext(c *Context) {
	s.Context = c
}

// 会话共同功能
type SessionAbility struct {
	val    any
	sendCh chan []byte
}

func (s *SessionAbility) Store(val any) {
	s.val = val
}

func (s *SessionAbility) Load() (val any, ok bool) {
	if s.val == nil {
		return
	}
	ok = true
	return
}

func (s *SessionAbility) InitSendChanel() {
	s.sendCh = make(chan []byte, 5)
}

func (s *SessionAbility) WriteSendChannel(buf []byte) {
	s.sendCh <- buf
}

func (s *SessionAbility) RunningSendLoop(handler func([]byte)) {
	go func() {
		for buf := range s.sendCh {
			handler(buf)
		}
	}()
}

func (s *SessionAbility) StopAbility() {
	close(s.sendCh)
	s.val = nil
}

// 核心会话标志
type SessionIdentify struct {
	*Context
	id        uint64
	closeFlag atomic.Bool
}

func (s *SessionIdentify) ClearIdentify() {
	s.Context = nil
	s.id = 0
	s.closeFlag.Store(false)
}

func (s *SessionIdentify) ID() uint64 {
	return s.id
}

func (s *SessionIdentify) SetID(id uint64) {
	s.id = id
}

func (s *SessionIdentify) UpdateID(id uint64) {
	value, ok := s.Context.sessionMgr.alive.Load(s.id)
	if ok {
		s.Context.sessionMgr.alive.Delete(s.id)
		s.id = id
		s.Context.sessionMgr.alive.Store(s.id, value)
	}
}

func (s *SessionIdentify) WithContext(c *Context) {
	s.Context = c
}

func (s *SessionIdentify) IsClosed() bool {
	return s.closeFlag.Load()
}

func (s *SessionIdentify) SetClosedStatus() {
	s.closeFlag.Store(true)
}
