package ws

import (
	"github.com/flylib/gonet"
	"github.com/gorilla/websocket"
	"net"
	"reflect"
)

var _ gonet.ISession = new(session)

// webSocket conn
type session struct {
	gonet.SessionCommon

	conn *websocket.Conn
	option
}

// 新会话
func newSession(c *gonet.Context, conn *websocket.Conn) *session {
	is := c.CreateSession()
	s := is.(*session)
	s.conn = conn
	s.WithContext(c)
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Close() error {
	return s.conn.Close()
}

// websocket does not support sending messages concurrently
func (s *session) Send(msgID uint32, msg any) (err error) {
	buf, err := s.Context.Package(s, msgID, msg)
	if err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	err = s.conn.WriteMessage(websocket.BinaryMessage, buf)
	return
}

// Loop to read messages
func (s *session) ReadLoop() {
	for {
		_, buf, err := s.conn.ReadMessage()
		if err != nil {
			s.Context.RecycleSession(s, err)
			return
		}
		msg, _, err := s.Context.UnPackage(s, buf)
		if err != nil {
			s.ILogger.Warnf("session_%v msg parser error,reason is %v ", s.ID(), err)
			continue
		}
		s.Context.PushGlobalMessageQueue(msg)
	}
}

func SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
