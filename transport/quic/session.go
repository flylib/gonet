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
func newSession(c *gonet.Context, conn quic.Connection) *session {
	ses := c.CreateSession()
	s, _ := ses.(*session)
	s.conn = conn
	s.WithContext(c)
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Send(msgID uint32, msg any) error {
	buf, err := s.Context.Package(s, msgID, msg)
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
			s.Warnf("session_%v AcceptStream error,reason is %v", s.ID(), err)
			s.Context.RecycleSession(s, err)
			return
		}
		s.channels = append(s.channels, ch)
		go s.recvLoop(ch)
	}

}

// 循环读取消息
func (s *session) recvLoop(channel quic.Stream) {
	defer channel.Close()
	var buf = make([]byte, gonet.MTU)
	for {
		n, err := channel.Read(buf)
		if err != nil {
			s.ILogger.Warnf("session_%v_%d steam read error - %v ", s.ID(), channel.StreamID(), err)
			channel.Close()
			return
		}
		msg, _, err := s.Context.UnPackage(s, buf[:n])
		if err != nil {
			s.ILogger.Warnf("session_%v_%d msg parser error,reason is %v ", s.ID(), channel.StreamID(), err)
			continue
		}
		msg.ID()

		s.Context.PushGlobalMessageQueue(msg)
	}
}

func SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
