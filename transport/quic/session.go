package quic

import (
	"context"
	"github.com/flylib/gonet"
	"github.com/quic-go/quic-go"
	"net"
)

// conn
type session struct {
	gonet.SessionIdentify
	gonet.SessionAbility
	conn   quic.Connection
	stream quic.Stream
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
	_, err = s.stream.Write(buf)
	return err
}

func (s *session) Close() error {
	return s.conn.CloseWithError(0, "EOF")
}

func (s *session) acceptStreamLoop() {
	for {
		channel, err := s.conn.AcceptStream(context.Background())
		if err != nil {
			s.Warnf("session_%v AcceptStream error,reason is %v", s.ID(), err)
			s.Context.RecycleSession(s, err)
			return
		}
		s.stream = channel
		go s.recvLoop(channel)
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
			return
		}
		msg, _, err := s.Context.UnPackage(s, buf[:n])
		if err != nil {
			s.ILogger.Warnf("session_%v_%d msg parser error,reason is %v ", s.ID(), channel.StreamID(), err)
			continue
		}
		s.Context.PushGlobalMessageQueue(msg)
	}
}
