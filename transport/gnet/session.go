package gnet

import (
	"github.com/flylib/gonet"
	"github.com/panjf2000/gnet/v2"
	"net"
	"reflect"
)

// Socket会话
type session struct {
	//核心会话标志
	gonet.SessionCommon
	//存储功能

	//累计收消息总数
	recvCount uint64
	//raw conn
	conn gnet.Conn
	//缓存数据，用于解决粘包问题
	cache []byte
}

// 新会话
func newSession(c *gonet.Context, conn gnet.Conn) *session {
	is := c.GetIdleSession()
	ns := is.(*session)
	ns.conn = conn
	ns.WithContext(c)
	ns.UpdateID(uint64(conn.Fd()))
	c.GetEventHandler().OnConnect(ns)
	return ns
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Send(msgID uint32, msg any) error {
	buf, err := s.GetContext().Package(s, msgID, msg)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(buf)
	return err
}

func (s *session) Close() error {
	return s.conn.Close()
}

func SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
