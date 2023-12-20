package kcp

import (
	"github.com/flylib/gonet"
	"github.com/xtaci/kcp-go"
	"net"
)

// Socket会话
type Session struct {
	//核心会话标志
	gonet.SessionCommon
	//存储功能

	//累计收消息总数
	recvCount uint64
	//raw conn
	conn *kcp.UDPSession
}

// 新会话
func newSession(c *gonet.Context, conn *kcp.UDPSession) *Session {
	is := c.CreateSession()
	s := is.(*Session)
	s.conn = conn
	s.WithContext(c)
	return s
}

func (s *Session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Session) Send(msgID uint32, msg any) error {
	buf, err := s.Context.Package(s, msgID, msg)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(buf)
	return err
}

func (s *Session) Close() error {
	return s.conn.Close()
}

// 接收循环
func (s *Session) recvLoop() {
	var buf = make([]byte, gonet.MTU)
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
