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
	SessionIdentify
	SessionStore
	conn *websocket.Conn
}

// 新会话
func newSession(c *Context, conn *websocket.Conn) *session {
	ses := c.CreateSession()
	newSession, _ := ses.(*session)
	newSession.conn = conn
	c.HandingMessage(newSession, &Message{
		ID: SessionConnect,
	})
	return newSession
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Close() error {
	return s.conn.Close()
}

// websocket does not support sending messages concurrently
func (s *session) Send(msg interface{}) error {
	return SendWSPacket(s.conn, msg)
}

// 循环读取消息
func (s *session) recvLoop(c *Context) {
	for {
		_, data, err := s.conn.ReadMessage()
		if err != nil {
			c.RecycleSession(s, err)
			return
		}
		msg, err := ParserWSPacket(data)
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		c.HandingMessage(s, msg)
	}
}
