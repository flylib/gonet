package ws

import (
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/transport"
	"net"
)

// webSocket conn
type session struct {
	SessionIdentify
	SessionStore
	conn *websocket.Conn
}

func init() {
	RegisterServer(&server{}, session{})
}

//新会话
func newSession(conn *websocket.Conn) *session {
	ses := CreateSession()
	s, _ := ses.(*session)
	s.conn = conn
	CacheMsg(&Message{
		Session: s,
		ID:      SessionConnect,
	})
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Close() error {
	err := s.conn.Close()
	s.conn = nil
	return err
}

//websocket does not support sending messages concurrently
func (s *session) Send(msg interface{}, params ...interface{}) error {
	return transport.SendWSPacket(s.conn, msg)
}

//循环读取消息
func (s *session) recvLoop() {
	for {
		_, pkt, err := s.conn.ReadMessage()
		if err != nil {
			RecycleSession(s, err)
			return
		}
		msg, err := transport.ParserWSPacket(pkt)
		if err != nil {
			CacheMsg(&Message{
				Session: s,
				ID:      SessionWarn,
				Body:    err,
			})
			continue
		}
		msg.Session = s
		CacheMsg(msg)
	}
}
