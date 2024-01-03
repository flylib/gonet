package ws

import (
	"github.com/flylib/gonet"
	"github.com/gorilla/websocket"
	"net"
	"reflect"
)

var _ gonet.ISession = new(session)

// webSocket conn
type session struct {
	gonet.SessionCommon

	conn *websocket.Conn
	option
}

// 新会话
func newSession(conn *websocket.Conn) *session {
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

func (s *session) Close() error {
	return s.conn.Close()
}

// websocket does not support sending messages concurrently
func (s *session) Send(msgID uint32, msg any) (err error) {
	buf, err := gonet.GetNetPackager().Package(msgID, msg)
	if err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	err = s.conn.WriteMessage(websocket.BinaryMessage, buf)
	return
}

// Loop to read messages
func (s *session) ReadLoop() {
	for {
		_, buf, err := s.conn.ReadMessage()
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

func SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
