package ws

import (
	"github.com/flylib/gonet"
	"github.com/gorilla/websocket"
	"net"
)

// Session is the WebSocket connection session.
type Session struct {
	gonet.SessionCommon
	conn *websocket.Conn
	option
}

// newSession gets an idle session from the pool and attaches conn.
// Returns nil if the max session limit has been reached.
func newSession(c *gonet.Context[*Session], conn *websocket.Conn) *Session {
	s, ok := c.GetIdleSession()
	if !ok {
		return nil
	}
	s.conn = conn
	c.GetEventHandler().OnConnect(s)
	return s
}

func (s *Session) RemoteAddr() net.Addr { return s.conn.RemoteAddr() }

func (s *Session) Close() error {
	return s.conn.Close()
}

// Send is safe to call from multiple goroutines (WebSocket requires serialized writes).
func (s *Session) Send(msgID uint32, msg any) error {
	buf, err := s.GetContext().Package(s, msgID, msg)
	if err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	return s.conn.WriteMessage(websocket.BinaryMessage, buf)
}

// Clear resets the session so it can be reused from the pool.
func (s *Session) Clear() {
	s.SessionCommon.Clear()
	s.conn = nil
}

// ReadLoop reads WebSocket frames until the connection closes.
func (s *Session) ReadLoop() {
	for {
		_, buf, err := s.conn.ReadMessage()
		if err != nil {
			s.GetContext().GetEventHandler().OnClose(s, err)
			s.GetContext().RecycleSession(s)
			return
		}
		msg, _, err := s.GetContext().UnPackage(s, buf)
		if err != nil {
			s.GetContext().GetEventHandler().OnError(s, err)
			continue
		}
		s.GetContext().PushGlobalMessageQueue(msg)
	}
}
