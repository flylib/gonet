package ws

import (
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"log"
	"net"
)

var _ ISession = new(session)

// webSocket conn
type session struct {
	*Context
	SessionIdentify
	SessionStore
	conn *websocket.Conn
}

// 新会话
func newSession(c *Context, conn *websocket.Conn) *session {
	ses := c.CreateSession()
	s := ses.(*session)
	s.conn = conn
	c.PushGlobalMessageQueue(s, NewSessionMessage)
	s.Context = c
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Close() error {
	s.Context = nil
	return s.conn.Close()
}

// websocket does not support sending messages concurrently
func (s *session) Send(msg any) error {
	bytes, err := s.Context.Package(msg)
	if err != nil {
		return err
	}
	return s.conn.WriteMessage(websocket.BinaryMessage, bytes)
}

// 循环读取消息
func (s *session) recvLoop() {
	for {
		_, data, err := s.conn.ReadMessage()
		if err != nil {
			s.Context.RecycleSession(s, err)
			return
		}
		//msg, err := ParserWSPacket(data)
		msg, err := s.Context.UnPackage(data)
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		s.Context.PushGlobalMessageQueue(s, msg)
	}
}
