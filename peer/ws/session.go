package ws

import (
	. "github.com/Quantumoffices/goNet"
	"github.com/Quantumoffices/goNet/codec"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

// webSocket session
type session struct {
	SessionIdentify
	SessionStore
	SessionActor
	socket *websocket.Conn
	sendCh chan interface{}
	closed bool
}

func init() {
	RegisterSessionType(session{})
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
		msgType, pkt, err := s.socket.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			logs.Warn("session_%d closed, %s", s.ID(), err)
			RecycleSession(s)
			//exit send goroutine
			s.sendCh <- nil
			break
		}
		actorID, msg, err := codec.ParserWSPacket(pkt)
		if err != nil {
			logs.Warn("msg parser error,reason is %v", err)
			continue
		}
		actor, err := s.GetActor(actorID)
		if err != nil {
			logs.Warn("session_%v get controller_%v error, reason is %v", s.ID(), actorID, err)
			continue
		}
		HandleEvent(NewEvent(EventNetWorkIO, s, actor, msg))
	}
}
