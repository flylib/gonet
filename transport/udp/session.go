package udp

import (
	"github.com/flylib/gonet"
	"net"
	"time"
)

// session is the UDP connection session.
type session struct {
	gonet.SessionCommon

	remoteAddr             *net.UDPAddr
	serverConn, remoteConn *net.UDPConn
	uuid                   string
	heartbeatTime          time.Time
	nexCheckTime           time.Time
}

// newSession gets an idle session from the pool and attaches the connection info.
func newSession(c *gonet.AppContext[*session], conn *net.UDPConn, remote *net.UDPAddr) *session {
	s, ok := c.GetIdleSession()
	if !ok {
		return nil
	}
	s.serverConn = conn
	s.remoteAddr = remote
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.remoteAddr
}

func (s *session) Send(msgID uint32, msg any) error {
	data, err := s.GetContext().Package(s, msgID, msg)
	if err != nil {
		return err
	}
	if s.remoteConn != nil {
		_, err = s.remoteConn.Write(data)
	} else {
		_, err = s.serverConn.WriteToUDP(data, s.remoteAddr)
	}
	return err
}

func (s *session) Close() error {
	return s.serverConn.Close()
}

// recvLoop reads from a dedicated UDP connection (client-side).
func (s *session) recvLoop() {
	var buf = make([]byte, 1024)
	for {
		n, err := s.serverConn.Read(buf)
		if err != nil {
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
