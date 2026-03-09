package gnet

import (
	"github.com/flylib/gonet"
	"github.com/panjf2000/gnet/v2"
	"net"
)

// session is the gnet connection session.
type session struct {
	gonet.SessionCommon

	recvCount uint64
	conn      gnet.Conn
	cache     []byte
}

// newSession gets an idle session from the pool and attaches the gnet connection.
func newSession(c *gonet.Context[*session], conn gnet.Conn) *session {
	s, ok := c.GetIdleSession()
	if !ok {
		return nil
	}
	s.conn = conn
	// Use the file descriptor as the session ID so we can look up sessions by conn.
	s.UpdateID(s, uint64(conn.Fd()))
	c.GetEventHandler().OnConnect(s)
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Send(msgID uint32, msg any) error {
	buf, err := s.GetContext().Package(s, msgID, msg)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(buf)
	return err
}

func (s *session) Close() error {
	return s.conn.Close()
}
