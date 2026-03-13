package kcp

import (
	"github.com/flylib/gonet"
	"github.com/xtaci/kcp-go"
	"net"
)

// Session is the KCP connection session.
type Session struct {
	gonet.SessionCommon
	recvCount uint64
	conn      *kcp.UDPSession
}

// newSession gets an idle session from the pool and attaches conn.
// Returns nil if the max session limit has been reached.
func newSession(c *gonet.AppContext[*Session], conn *kcp.UDPSession) *Session {
	s, ok := c.GetIdleSession()
	if !ok {
		return nil
	}
	s.conn = conn
	c.GetEventHandler().OnConnect(s)
	return s
}

func (s *Session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Session) Send(msgID uint32, msg any) error {
	buf, err := s.GetContext().Package(s, msgID, msg)
	if err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	_, err = s.conn.Write(buf)
	return err
}

func (s *Session) Close() error {
	return s.conn.Close()
}

// Clear resets the session so it can be reused from the pool.
func (s *Session) Clear() {
	s.SessionCommon.Clear()
	s.conn = nil
	s.recvCount = 0
}

// recvLoop reads from the KCP connection until it closes.
func (s *Session) recvLoop() {
	buf := make([]byte, gonet.MTU)
	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			s.GetContext().GetEventHandler().OnClose(s, err)
			s.GetContext().RecycleSession(s)
			return
		}
		msg, _, err := s.GetContext().UnPackage(s, buf[:n])
		if err != nil {
			s.GetContext().GetEventHandler().OnError(s, err)
			continue
		}
		s.GetContext().PushGlobalMessageQueue(msg)
	}
}
