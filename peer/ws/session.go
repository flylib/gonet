package ws

import (
	. "github.com/Quantumoffices/goNet"
	"github.com/Quantumoffices/goNet/codec"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

// webSocket session
type session struct {
	SessionIdentify //标志
	SessionStore    //存储
	*websocket.Conn //socket
	sendCh          chan interface{}
	closed          bool //关闭标志
}

func init() {
	RegisterSessionType(session{})
}

func newSession(conn *websocket.Conn) *session {
	newSession := AddSession().(*session)
	newSession.Conn = conn
	newSession.sendCh = make(chan interface{}, 1)
	return newSession
}

func (s *session) Socket() interface{} {
	return s.Conn
}

func (s *session) Close() {
	if s.closed {
		return
	}
	s.closed = true
	if err := s.Conn.Close(); err != nil {
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
		if err := codec.SendWSPacket(s.Conn, msg); err != nil {
			logs.Error("sesssion_%v send msg error,reason is %v", s.ID(), err)
			break
		}
	}
}

//read
func (s *session) recvLoop() {
	for {
		msgType, pkt, err := s.Conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			logs.Warn("session_%d closed, %s", s.ID(), err)
			RecycleSession(s)
			//exit send goroutine
			s.sendCh <- nil
			break
		}
		msg, err := codec.ParserWSPacket(pkt)
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
