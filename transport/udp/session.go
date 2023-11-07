package udp

import (
	"github.com/flylib/gonet"
	"net"
	"reflect"
)

//addr:sessionID
var remotes = map[string]uint64{}

// Socket会话
type session struct {
	gonet.SessionIdentify
	gonet.SessionAbility
	remoteAddr       *net.UDPAddr
	conn, remoteConn *net.UDPConn
	data             interface{}
	buf              []byte
}

// 新会话
func newSession(c *gonet.Context, conn *net.UDPConn, remote *net.UDPAddr) *session {
	is := c.CreateSession()
	s := is.(*session)
	s.conn = conn
	s.remoteAddr = remote
	s.buf = make([]byte, gonet.MTU)
	s.WithContext(c)
	remotes[remote.String()] = s.ID()
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.remoteAddr
}

// 发送封包
func (s *session) Send(msgID uint32, msg any) error {
	data, err := s.Context.Package(s, msgID, msg)
	if err != nil {
		return err
	}
	if s.remoteConn != nil {
		_, err = s.remoteConn.Write(data)
	} else {
		_, err = s.conn.WriteToUDP(data, s.remoteAddr)
	}

	return err
}

func (s *session) Close() error {
	return s.conn.Close()
}

// Loop to read messages
func (s *session) recvLoop() {
	var buf = make([]byte, 1024)
	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			s.Context.RecycleSession(s, err)
			return
		}
		msg, _, err := s.Context.UnPackage(s, buf[:n])
		if err != nil {
			s.ILogger.Warnf("session_%v msg parser error,reason is %v ", s.ID(), err)
			continue
		}
		s.Context.PushGlobalMessageQueue(msg)
	}
}

// todo 心跳检测
func SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
