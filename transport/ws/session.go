package ws

import (
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/codec"
	"log"
)

// webSocket conn
type session struct {
	SessionIdentify
	SessionStore
	conn *websocket.Conn
}

func init() {
	SetSessionType(&session{})
}

//新会话
func newSession(conn *websocket.Conn) *session {
	ses := CreateSession()
	newSession, _ := ses.(*session)
	newSession.conn = conn
	return newSession
}

func (s *session) Close() error {
	return s.conn.Close()
}

//websocket does not support sending messages concurrently
func (s *session) Send(msg interface{}) error {
	return codec.SendWSPacket(s.conn, msg)
}

//循环读取消息
func (s *session) recvLoop() {
	for {
		_, pkt, err := s.conn.ReadMessage()
		if err != nil {
			log.Printf("session_%v closed, %v \n", s.ID(), err)
			RecycleSession(s)
			return
		}
		msg, err := codec.ParserWSPacket(pkt)
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		msg.Session = s
		PushWorkerPool(msg)
	}
}
