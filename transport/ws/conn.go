package ws

import (
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/transport"
	"log"
	"net"
)

// webSocket conn
type Conn struct {
	ConnIdentify
	ConnStore
	conn *websocket.Conn
}

func init() {
	RegisterServer(&server{}, Conn{})
}

//新会话
func newConn(conn *websocket.Conn) *Conn {
	c := CreateConn()
	newConnection, _ := c.(*Conn)
	newConnection.conn = conn
	CacheSession(&Session{
		Connection: newConnection,
		Msg:        GetCommonMsgNewConnMsg(),
	})
	return newConnection
}

func (s *Conn) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Conn) Close() error {
	return s.conn.Close()
}

//websocket does not support sending messages concurrently
func (s *Conn) Send(msg interface{}) error {
	return transport.SendWSPacket(s.conn, msg)
}

//循环读取消息
func (s *Conn) recvLoop() {
	for {
		_, pkt, err := s.conn.ReadMessage()
		if err != nil {
			RecycleConn(s, err)
			return
		}
		msg, err := transport.ParserWSPacket(pkt)
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		CacheSession(&Session{
			Connection: s,
			Msg:        msg,
		})
	}
}
