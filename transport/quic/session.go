package quic

import (
	"context"
	"fmt"
	"github.com/quic-go/quic-go"
	. "github.com/zjllib/gonet/v3"
	"log"
	"net"
	"sync"
)

// webSocket conn
type session struct {
	SessionIdentify
	SessionStore
	conn    quic.Connection
	streams sync.Map
}

// 新会话
func newSession(c *Context, conn quic.Connection) *session {
	ses := c.CreateSession()
	s, _ := ses.(*session)
	s.conn = conn
	s.WithContext(c)
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Send(msg any) error {
	data, err := s.Context.Package(msg)
	if err != nil {
		return err
	}
	s.streams.Range(func(key, value any) bool {
		_, err = value.(quic.Stream).Write(data)
		if err != nil {
			return false
		}
		return true
	})
	//_, err = s.conn.Write(data)
	return err
}

func (s *session) Close() error {
	err := s.conn.CloseWithError(0, "EOF")
	s.conn = nil
	return err
}

// 循环读取消息
func (s *session) readStreamLoop(stream quic.Stream) {
	for {
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			log.Printf("session_%v stream_%v reading error,reason is %v \n", s.ID(), stream.StreamID(), err)
			err = stream.Close()
			if err != nil {
				log.Printf("session_%v stream_%v close error,reason is %v \n", s.ID(), stream.StreamID(), err)
			}
			s.streams.Delete(stream.StreamID())
			return
		}
		msg, _, err := s.Context.UnPackage(buf[:n])
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		s.Context.PushGlobalMessageQueue(s, msg)
	}
}

// stream
func (s *session) acceptStream() {
	for {
		stream, err := s.conn.AcceptStream(context.Background())
		if err != nil {
			s.Context.RecycleSession(s, err)
			return
		}
		fmt.Println(stream.StreamID())
		s.Store(stream.StreamID(), stream)
		//开单独go routine 去处理stream的消息
		go s.readStreamLoop(stream)
	}
}
