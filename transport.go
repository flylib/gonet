package gonet

import (
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

var (
	_ ISessionIdentify = new(SessionIdentify)
	_ ISessionAbility  = new(SessionAbility)
	_ IPeerIdentify    = new(PeerIdentify)
)

// 端属性
type PeerIdentify struct {
	*AppContext
	uuid string
	//地址
	addr string
}

func (s *PeerIdentify) Addr() string {
	return s.addr
}

func (s *PeerIdentify) SetAddr(addr string) {
	s.addr = addr
}

func (s *PeerIdentify) WithContext(c *AppContext) {
	s.AppContext = c
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
	s.val = nil
}

func (s *SessionAbility) PushSendChannel(buf []byte) {
	s.sendCh <- buf
}

func (s *SessionAbility) SendLoop(writeDataHandler func([]byte)) {
	for buf := range s.sendCh {
		writeDataHandler(buf)
	}
}

func (s *SessionAbility) StopAbility() {
	close(s.sendCh)
	s.val = nil
}

// 核心会话标志
type SessionIdentify struct {
	*AppContext
	id        uint64
	closeFlag atomic.Bool
}

func (s *SessionIdentify) ClearIdentify() {
	s.AppContext = nil
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
	value, ok := s.AppContext.sessionMgr.alive.Load(s.id)
	if ok {
		s.AppContext.sessionMgr.alive.Delete(s.id)
		s.id = id
		s.AppContext.sessionMgr.alive.Store(s.id, value)
	}
}

func (s *SessionIdentify) WithContext(c *AppContext) {
	s.AppContext = c
}

func (s *SessionIdentify) IsClosed() bool {
	return s.closeFlag.Load()
}

func (s *SessionIdentify) SetClosedStatus() {
	s.closeFlag.Store(true)
}
