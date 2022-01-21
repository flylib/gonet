package quic

import (
	"context"
	"github.com/lucas-clemente/quic-go"
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/transport"
	"log"
	"net"
)

// webSocket conn
type session struct {
	SessionIdentify
	SessionStore
	conn    quic.Connection
	streams map[quic.StreamID]quic.Stream
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
	} else {
		streamID, ok := params[0].(quic.StreamID)
		if ok {
			stream := s.streams[streamID]
			if stream != nil {
				err = transport.SendPacket(stream, msg)
			}
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
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		msg.Session = s
		CacheMsg(msg)
	}
}

//循环读取消息
func (s *session) recvStreamMsgLoop(stream quic.Stream) {
	for {
		var buf []byte
		n, err := stream.Read(buf)
		if err != nil {
			//RecycleSession(s, err)
			stream.Close()
			delete(s.streams, stream.StreamID())
			return
		}
		msg, _, err := transport.ParserTcpPacket(buf[:n])
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		msg.Session = s
		CacheMsg(msg)
	}
}

//
func (s *session) recvStreamLoop() {
	for {
		stream, err := s.conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("[quic] session_%v recvStreamLoop error,reason is %v \n", s.ID(), err)
			s.Close()
			break
		}
		s.streams[stream.StreamID()] = stream
	}
}
