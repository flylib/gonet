package udp

import (
	"github.com/flylib/gonet"
	"net"
	"reflect"
	"time"
)

// Socket会话
type session struct {
	gonet.SessionCommon

	remoteAddr             *net.UDPAddr
	serverConn, remoteConn *net.UDPConn
	uuid                   string
	heartbeatTime          time.Time //最近心跳时间点
	nexCheckTime           time.Time //下次检查时间点
}

// 新会话
func newSession(conn *net.UDPConn, remote *net.UDPAddr) *session {
	is := gonet.GetSessionManager().GetIdleSession()
	s := is.(*session)
	s.serverConn = conn
	s.remoteAddr = remote
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.remoteAddr
}

// 发送封包
func (s *session) Send(msgID uint32, msg any) error {
	data, err := gonet.GetNetPackager().Package(msgID, msg)
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

// Loop to read messages
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

// todo 心跳检测
func SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
