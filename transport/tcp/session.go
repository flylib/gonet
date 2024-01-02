package tcp

import (
	"encoding/binary"
	"github.com/flylib/gonet"
	"io"
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
	conn net.Conn
}

// 新会话
func newSession(conn net.Conn) *session {
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
func (s *session) readLoop() {
	for {
		buf, err := s.read()
		if err != nil {
			gonet.GetEventHandler().OnClose(s, err)
			gonet.GetSessionManager().RecycleSession(s)
			return
		}
		msg, err := gonet.GetNetPackager().UnPackage(s, buf)
		if err != nil {
			gonet.GetEventHandler().OnError(s, err)
			continue
		}
		gonet.GetAsyncRuntime().PushMessage(msg)
	}
}

// 粘包处理
func (s *session) read() ([]byte, error) {
	header := make([]byte, gonet.PktSizeLen)
	_, err := io.ReadFull(s.conn, header)
	if err != nil {
		return nil, err
	}

	bodyLength := binary.LittleEndian.Uint16(header)
	buf := make([]byte, bodyLength)
	_, err = io.ReadFull(s.conn, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
