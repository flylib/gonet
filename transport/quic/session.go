package quic

import (
	"context"
	"github.com/lucas-clemente/quic-go"
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/transport"
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

func init() {
	RegisterServer(&server{}, session{})
}

//新会话
func newSession(conn quic.Connection) *session {
	ses := CreateSession()
	s, _ := ses.(*session)
	s.conn = conn
	CacheMsg(&Message{
		Session: s,
		ID:      SessionConnect,
	})
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Close() error {
	err := s.conn.CloseWithError(0, "EOF")
	s.conn = nil
	return err
}

//websocket does not support sending messages concurrently
func (s *session) Send(msg interface{}, params ...interface{}) error {
	var err error
	if len(params) == 0 {
		err = transport.SendQUICPacket(s.conn, msg)
	} else if streamID, ok := params[0].(quic.StreamID); ok {
		if load, ok := s.streams.Load(streamID); ok {
			err = transport.SendPacket(load.(quic.Stream), msg)
		}
	}
	return err
}

//循环读取消息
func (s *session) recvLoop() {
	for {
		bytes, err := s.conn.ReceiveMessage()
		if err != nil {
			RecycleSession(s, err)
			return
		}
		msg, _, err := transport.ParserTcpPacket(bytes)
		if err != nil {
			log.Printf("[quic]session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		msg.Session = s
		CacheMsg(msg)
	}
}

//循环读取消息
func (s *session) recvStreamMsgLoop(stream quic.Stream) {
	for {
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			CacheMsg(&Message{
				Session:  s,
				StreamID: stream.StreamID(),
				ID:       SessionWarn,
				Body:     err,
			})
			stream.Close()
			s.streams.Delete(stream.StreamID())
			return
		}
		msg, _, err := transport.ParserTcpPacket(buf[:n])
		if err != nil {
			CacheMsg(&Message{
				Session:  s,
				StreamID: stream.StreamID(),
				ID:       SessionWarn,
				Body:     err,
			})
			continue
		}
		msg.StreamID = stream.StreamID()
		msg.Session = s
		CacheMsg(msg)
	}
}

//stream
func (s *session) recvStreamLoop() {
	for {
		stream, err := s.conn.AcceptStream(context.Background())
		if err != nil {
			RecycleSession(s, err)
			return
		}
		s.Store(stream.StreamID(), stream)
		//开单独go routine 去处理stream的消息
		go s.recvStreamMsgLoop(stream)
	}
}
