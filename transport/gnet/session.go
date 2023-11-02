package gnet

import (
	"github.com/flylib/gonet"
	"github.com/panjf2000/gnet/v2"
	"net"
)

// Socket会话
type Session struct {
	//核心会话标志
	gonet.SessionIdentify
	//存储功能
	gonet.SessionAbility
	//累计收消息总数
	recvCount uint64
	//raw conn
	conn gnet.Conn
	//缓存数据，用于解决粘包问题
	cache []byte
}

// 新会话
func newSession(ctx *gonet.Context, conn gnet.Conn) *Session {
	is := ctx.CreateSession()
	s := is.(*Session)
	s.conn = conn
	s.WithContext(ctx)
	s.UpdateID(uint64(conn.Fd()))
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
