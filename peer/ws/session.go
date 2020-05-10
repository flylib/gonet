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
	SessionController
	conn *websocket.Conn
	buf  []byte
	sync.RWMutex
}

func newSession(conn *websocket.Conn) *session {
	ses := SessionManager.GetIdleSession()
	if ses == nil {
		ses = &session{
			buf: make([]byte, codec.MTU),
		}
		SessionManager.AddSession(ses)
	}
	ses.(*session).conn = conn
	return ses.(*session)
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
	//传入消息校验
	if msg == nil {
		return
	}
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
			SessionManager.RecycleSession(s)
			break
		}
		controllerIdx, msg, err := codec.ParserWSPacket(data)
		if err != nil {
			logs.Warn("msg parser error,reason is %v", err)
			continue
		}
		controller, err := s.GetController(controllerIdx)
		if err != nil {
			logs.Warn("session_%v get controller_%v error, reason is %v", s.ID(), controllerIdx, err)
			continue
		}
		SubmitMsgToAntsPool(controller, s, msg)
	}
}
