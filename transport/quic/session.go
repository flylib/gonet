package quic

import (
	"context"
	"fmt"
	"github.com/flylib/gonet"
	"github.com/quic-go/quic-go"
	"net"
)

// session is the QUIC connection session.
type session struct {
	gonet.SessionCommon

	conn     quic.Connection
	channels []quic.Stream
	mod      uint32
}

// newSession gets an idle session from the pool and attaches the QUIC connection.
func newSession(c *gonet.AppContext[*session], conn quic.Connection) *session {
	s, ok := c.GetIdleSession()
	if !ok {
		return nil
	}
	s.conn = conn
	c.GetEventHandler().OnConnect(s)
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Send(msgID uint32, msg any) error {
	buf, err := s.GetContext().Package(s, msgID, msg)
	if err != nil {
		return err
	}
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
			s.GetContext().GetEventHandler().OnClose(s, err)
			s.GetContext().RecycleSession(s)
			return
		}
		s.channels = append(s.channels, ch)
		go s.recvLoop(ch)
	}
}

// recvLoop reads messages from a single QUIC stream.
func (s *session) recvLoop(channel quic.Stream) {
	defer channel.Close()
	var buf = make([]byte, gonet.MTU)
	for {
		n, err := channel.Read(buf)
		if err != nil {
			s.GetContext().GetEventHandler().OnError(s, err)
			channel.Close()
			return
		}
		msg, _, err := s.GetContext().UnPackage(s, buf[:n])
		if err != nil {
			s.GetContext().GetEventHandler().OnError(s, err)
			continue
		}
		_ = msg.ID()
		s.GetContext().PushGlobalMessageQueue(msg)
	}
}
