package gnet

import (
	"github.com/panjf2000/gnet/v2"
	. "github.com/zjllib/gonet/v3"
	"net"
)

var _ ISession = new(session)

// Socket会话
type session struct {
	*Context
	//核心会话标志
	SessionIdentify
	//存储功能
	SessionAbility
	//累计收消息总数
	recvCount uint64
	//raw conn
	conn gnet.Conn
	//缓存数据，用于解决粘包问题
	cache []byte
}

// 新会话
func newSession(c *Context, conn gnet.Conn) *session {
	is := c.CreateSession()
	s := is.(*session)
	s.conn = conn
	s.WithContext(c)
	s.UpdateID(uint64(conn.Fd()))
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Send(msg interface{}) error {
	bytes, err := s.Context.Package(msg)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(bytes)
	return err
}

func (s *session) Close() error {
	return s.conn.Close()
}
