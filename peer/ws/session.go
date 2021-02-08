package ws

import (
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	. "github.com/zjllib/goNet"
	"github.com/zjllib/goNet/codec"
)

// webSocket session
type session struct {
	SessionIdentify //标志
	SessionStore    //存储
	SessionScene
	socket *websocket.Conn //socket
	sendCh chan interface{}
	closed bool //关闭标志
}

func init() {
	SetSessionType(session{})
}

func newSession(conn *websocket.Conn) *session {
	newSession := AddSession().(*session)
	newSession.socket = conn
	newSession.sendCh = make(chan interface{}, 1)
	return newSession
}

func (s *session) Socket() interface{} {
	return s.socket
}

func (s *session) Close() {
	if s.closed {
		return
	}
	s.closed = true
	if err := s.socket.Close(); err != nil {
		logs.Error("sesssion_%v close error,reason is %v", s.ID(), err)
	}
}

//websocket does not support sending messages concurrently
func (s *session) Send(msg interface{}) {
	//sending empty messages is not allowed
	if !s.closed && msg == nil {
		return
	}
	s.sendCh <- msg
}

//write
func (s *session) sendLoop() {
	for msg := range s.sendCh {
		if msg == nil {
			break
		}
		if err := codec.SendWSPacket(s.socket, msg); err != nil {
			logs.Error("sesssion_%v send msg error,reason is %v", s.ID(), err)
			break
		}
	}
}

//read
func (s *session) recvLoop() {
	for {
		wsMsgKind, pkt, err := s.socket.ReadMessage()
		if err != nil || wsMsgKind == websocket.CloseMessage {
			logs.Warn("session_%d closed, %s", s.ID(), err)
			RecycleSession(s)
			s.sendCh <- nil
			break
		}
		msg, err := codec.ParserWSPacket(pkt)
		if err != nil {
			logs.Warn("msg parser error,reason is %v", err)
			continue
		}
		msg.Session = s
		PushWorkerPool(msg)
	}
}
