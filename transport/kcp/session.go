package kcp

import (
	"github.com/xtaci/kcp-go/v5"
	. "github.com/zjllib/gonet/v3"
	"log"
	"net"
)

var _ ISession = new(session)

// Socket会话
type session struct {
	//核心会话标志
	SessionIdentify
	//存储功能
	SessionAbility
	//累计收消息总数
	recvCount uint64
	//raw conn
	conn *kcp.UDPSession
}

// 新会话
func newSession(c *AppContext, conn *kcp.UDPSession) *session {
	is := c.CreateSession()
	s := is.(*session)
	s.conn = conn
	s.WithContext(c)
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Send(msg interface{}) error {
	data, err := s.AppContext.Package(msg)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(data)
	return err
}

func (s *session) Close() error {
	return s.conn.Close()
}

// 接收循环
func (s *session) recvLoop() {
	for {
		var buf = make([]byte, MTU)
		n, err := s.conn.Read(buf)
		if err != nil {
			s.AppContext.RecycleSession(s, err)
			return
		}
		if n == 0 {
			continue
		}
		msg, _, err := s.UnPackage(buf[:n])
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		s.AppContext.PushGlobalMessageQueue(s, msg)
	}
}
