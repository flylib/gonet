package udp

import (
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/transport"
	"net"
)

var _ transport.ISession = new(session)

//addr:sessionID
var remotes = map[string]uint64{}

// Socket会话
type session struct {
	transport.SessionIdentify
	transport.SessionStore
	remote *net.UDPAddr
	conn   *net.UDPConn
	data   interface{}
	buf    []byte
}

//新会话
func newSession(conn *net.UDPConn, remote *net.UDPAddr) *session {
	ses := CreateSession()
	ses.(*session).conn = conn
	remotes[remote.String()] = ses.ID()
	return ses.(*session)
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

// 发送封包
func (s *session) Send(msg interface{}) error {
	var err error
	if s.remote == nil {
		err = SendPacket(s.conn, msg)
	} else {
		err = SendUdpPacket(s.conn, msg, s.remote)
	}
	return err
}

func (s *session) Close() error {
	return s.conn.Close()
}
