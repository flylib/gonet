package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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
		logrus.Errorf("sesssion_%v close error,reason is %v", s.ID(), err)
	}
	s.data = nil
}

// 发送封包
func (s *session) Send(msg interface{}) {
	s.Lock()
	defer s.Unlock()
	if err := codec.SendWSPacket(s.conn, msg); err != nil {
		logrus.Errorf("sesssion_%v send msg error,reason is %v", s.ID(), err)
		logrus.Errorf(s.conn.RemoteAddr().String())
	}
}

// 接收循环
func (s *session) recvLoop() {
	for {
		t, body, err := s.conn.ReadMessage()
		if err != nil || t == websocket.CloseMessage {
			logrus.Warnf("session_%d closed, err: %s", s.ID(), err)
			goNet.SessionManager.RecycleSession(s)
			break
		}
		var msg goNet.Message
		switch t {
		case websocket.TextMessage:
			//logrus.Info("TextMessage")
			msg, err = codec.ParserWSPacket(body)
			if err != nil {
				logrus.Warnf("message decode error=%v", err)
				continue
			}
		case websocket.BinaryMessage:
			//logrus.Info("BinaryMessage")
			msg, err = codec.ParserPacket(body)
			if err != nil {
				logrus.Warnf("message decode error=%s", err)
				continue
			}
		default:
			logrus.Errorf("unknown message")
			continue
		}
		goNet.HandleMessage(msg, s)
	}
}

func (u *session) Value(v ...interface{}) interface{} {
	if len(v) > 0 {
		u.data = v[0]
	}
	return u.data
}
