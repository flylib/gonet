package ws

import (
	"github.com/gorilla/websocket"
	"goNet"
	"goNet/codec"
	"sync"
)

// webSocket session
type session struct {
	goNet.SessionIdentify
	//core connection
	conn *websocket.Conn
	data interface{}
	buf  []byte
	sync.RWMutex
	//example center_service/room_service/...
	//stubs []interface{}
}

func newSession(conn *websocket.Conn) *session {
	ses := goNet.SessionManager.GetIdleSession()
	if ses == nil {
		ses = &session{
			buf: make([]byte, codec.MTU),
		}
		goNet.SessionManager.AddSession(ses)
	}
	ses.(*session).conn = conn
	return ses.(*session)
}

// 取原始连接
func (s *session) Socket() interface{} {
	return s.conn
}

func (s *session) Close() {
	if err := s.conn.Close(); err != nil {
		goNet.Log.Errorf("sesssion_%v close error,reason is %v", s.ID(), err)
	}
	s.data = nil
}

func (s *session) Send(msg interface{}) {
	//传入消息校验
	if msg == nil {
		return
	}
	s.Lock()
	defer s.Unlock()
	if err := codec.SendWSPacket(s.conn, msg); err != nil {
		goNet.Log.Errorf("sesssion_%v send msg error,reason is %v", s.ID(), err)
		goNet.Log.Errorf(s.conn.RemoteAddr().String())
	}
}

// 接收循环
func (s *session) recvLoop() {
	for {
		t, data, err := s.conn.ReadMessage()
		if err != nil || t == websocket.CloseMessage {
			goNet.Log.Warnf("session_%d closed, err: %s", s.ID(), err)
			goNet.SessionManager.RecycleSession(s)
			break
		}
		var msg goNet.Msg
		msg, err = codec.ParserWSPacket(data)
		if err != nil {
			goNet.Log.Warnf("parse message error:%s", err)
			continue
		}
		goNet.SubmitMsgToAntsPool(msg, s)
	}
}

func (u *session) Value(v ...interface{}) interface{} {
	if len(v) > 0 {
		u.data = v[0]
	}
	return u.data
}
