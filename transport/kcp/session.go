package kcp

import (
	"github.com/flylib/gonet"
	"github.com/xtaci/kcp-go"
	"net"
)

// Socket会话
type session struct {
	//核心会话标志
	gonet.SessionCommon
	//存储功能

	//累计收消息总数
	recvCount uint64
	//raw conn
	conn *kcp.UDPSession
}

// 新会话
func newSession(conn *kcp.UDPSession) *session {
	is := gonet.GetSessionManager().GetIdleSession()
	ns := is.(*session)
	ns.conn = conn
	gonet.GetSessionManager().AddSession(ns)
	gonet.GetEventHandler().OnConnect(ns)
	return ns
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Send(msgID uint32, msg any) error {
	buf, err := gonet.GetNetPackager().Package(msgID, msg)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(buf)
	return err
}

func (s *session) Close() error {
	return s.conn.Close()
}

// 接收循环
func (s *session) recvLoop() {
	var buf = make([]byte, gonet.MTU)
	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			gonet.GetEventHandler().OnClose(s, err)
			gonet.GetSessionManager().RecycleSession(s)
			return
		}
		msg, err := gonet.GetNetPackager().UnPackage(s, buf[:n])
		if err != nil {
			gonet.GetEventHandler().OnError(s, err)
			continue
		}
		gonet.GetAsyncRuntime().PushMessage(msg)
	}
}
