package ws

import (
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/codec"
	"sync"
)

// webSocket conn
type conn struct {
	SessionIdentify //标志
	sync.Map        //存储
	SessionScene
	socket *websocket.Conn //socket
	sendCh chan interface{}
	closed bool //关闭标志
}

func init() {
	SetconnType(conn{})
}

func newconn(conn *websocket.Conn) *conn {
	newconn := Addconn().(*conn)
	newconn.socket = conn
	newconn.sendCh = make(chan interface{}, 1)
	return newconn
}

func (s *conn) Socket() interface{} {
	return s.socket
}

func (s *conn) Close() {
	if s.closed {
		return
	}
	s.closed = true
	if err := s.socket.Close(); err != nil {
		logs.Error("sesssion_%v close error,reason is %v", s.ID(), err)
	}
}

//websocket does not support sending messages concurrently
func (s *conn) Send(msg interface{}) {
	//sending empty messages is not allowed
	if !s.closed && msg == nil {
		return
	}
	s.sendCh <- msg
}

//write
func (s *conn) sendLoop() {
	for msg := range s.sendCh {
		if msg == nil {
			break
		}
		if err := codec.SendWSPacket(s.socket, msg); err != nil {
			logs.Error("sesssion_%v send msg error,reason is %v", s.ID(), err)
			break
		}
	}
}

//read
func (s *conn) recvLoop() {
	for {
		wsMsgKind, pkt, err := s.socket.ReadMessage()
		if err != nil || wsMsgKind == websocket.CloseMessage {
			logs.Warn("conn_%d closed, %s", s.ID(), err)
			Recycleconn(s)
			s.sendCh <- nil
			break
		}
		msg, err := codec.ParserWSPacket(pkt)
		if err != nil {
			logs.Warn("msg parser error,reason is %v", err)
			continue
		}
		msg.conn = s
		PushWorkerPool(msg)
	}
}
