package udp

import (
	. "github.com/zjllib/gonet/v3"
	"net"
)

var _ ISession = new(session)

//addr:sessionID
var remotes = map[string]uint64{}

// Socket会话
type session struct {
	SessionIdentify
	SessionAbility
	remoteAddr *net.UDPAddr
	conn       *net.UDPConn
	data       interface{}
	buf        []byte
}

// 新会话
func newSession(c *Context, conn *net.UDPConn, remote *net.UDPAddr) *session {
	is := c.CreateSession()
	s := is.(*session)
	s.conn = conn
	s.remoteAddr = remote
	s.buf = make([]byte, MTU)
	s.WithContext(c)
	remotes[remote.String()] = s.ID()
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

// 发送封包
func (s *session) Send(msg interface{}) error {
	data, err := s.Context.Package(msg)
	if err != nil {
		return err
	}
	_, err = s.conn.WriteToUDP(data, s.remoteAddr)
	return err
}

func (s *session) Close() error {
	return s.conn.Close()
}

//todo 心跳检测
