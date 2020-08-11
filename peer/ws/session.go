package ws

import (
	. "github.com/Quantumoffices/goNet"
	"github.com/Quantumoffices/goNet/codec"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"sync"
)

// webSocket session
type session struct {
	SessionIdentify
	SessionStore
	SessionRoute
	conn *websocket.Conn
	sync.RWMutex
}

func init() {
	RegisterSessionType(session{})
}

func newSession(conn *websocket.Conn) *session {
	newSession := AddSession().(*session)
	newSession.conn = conn
	return newSession
}

func (s *session) Socket() interface{} {
	return s.conn
}

func (s *session) Close() {
	if err := s.conn.Close(); err != nil {
		logs.Error("sesssion_%v close error,reason is %v", s.ID(), err)
	}
}

func (s *session) Send(msg interface{}) {
	s.Lock()
	defer s.Unlock()
	if err := codec.SendWSPacket(s.conn, msg); err != nil {
		logs.Error("sesssion_%v send msg error,reason is %v", s.ID(), err)
	}
}

// 接收循环
func (s *session) recvLoop() {
	for {
		t, data, err := s.conn.ReadMessage()
		if err != nil || t == websocket.CloseMessage {
			logs.Warn("session_%d closed, err: %s", s.ID(), err)
			RecycleSession(s)
			break
		}
		routeID, msg, err := codec.ParserWSPacket(data)
		if err != nil {
			logs.Warn("msg parser error,reason is %v", err)
			continue
		}
		route, err := s.GetRoute(routeID)
		if err != nil {
			logs.Warn("session_%v get controller_%v error, reason is %v", s.ID(), routeID, err)
			continue
		}
		//HandleEvent(NewEvent(s, controller, msg))
		CommitWorkerPool(Event{From: s, Router: route, Msg: msg})
	}
}
