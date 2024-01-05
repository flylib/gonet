package quic

import (
	"context"
	"fmt"
	"github.com/flylib/gonet"
	"github.com/quic-go/quic-go"
	"net"
	"reflect"
)

// conn
type session struct {
	gonet.SessionCommon

	conn     quic.Connection
	channels []quic.Stream
	mod      uint32
}

// 新会话
func newSession(conn quic.Connection) *session {
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
	//2*1000000+101 = 2000000101
	channelId := msgID / s.mod
	msgID = msgID - channelId*s.mod
	if s.channels[channelId] != nil {
		_, err = s.channels[channelId].Write(buf)
		if err != nil {
			return err
		}
	}

	return fmt.Errorf("not found the channel-%d", channelId)
}

func (s *session) Close() error {
	return s.conn.CloseWithError(0, "EOF")
}

func (s *session) acceptStream() {
	for {
		ch, err := s.conn.AcceptStream(context.Background())
		if err != nil {
			gonet.GetEventHandler().OnClose(s, err)
			gonet.GetSessionManager().RecycleSession(s)
			return
		}
		s.channels = append(s.channels, ch)
		go s.readLoop(ch)
	}

}

// 循环读取消息
func (s *session) readLoop(channel quic.Stream) {
	defer channel.Close()
	var buf = make([]byte, gonet.MTU)
	for {
		n, err := channel.Read(buf)
		if err != nil {
			gonet.GetEventHandler().OnError(s, err)
			channel.Close()
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

func SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
