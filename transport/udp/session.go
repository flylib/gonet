package udp

import (
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/transport"
	"net"
)

//addr:sessionID
var remotes = map[string]uint64{}

// Socket会话
type session struct {
	SessionIdentify
	SessionStore
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
		err = transport.SendPacket(s.conn, msg)
	} else {
		err = transport.SendUdpPacket(s.conn, msg, s.remote)
	}
	return err
}

func (s *session) Close() error {
	return s.conn.Close()
}
