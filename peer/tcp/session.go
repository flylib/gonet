package tcp

import (
	"github.com/sirupsen/logrus"
	"goNet"
	"goNet/codec"
	"net"
)

// Socket会话
type session struct {
	//核心会话标志
	goNet.SessionIdentify
	goNet.SessionController
	//累计收消息总数
	recvCount uint64
	//raw conn
	conn  net.Conn
	store interface{}
	buf   []byte
	//缓存数据，用于解决粘包问题
	//cache []byte
}

//新会话
func newSession(conn net.Conn) *session {
	ses := goNet.SessionManager.GetIdleSession()
	if ses == nil {
		ses = &session{
			conn: conn,
			buf:  make([]byte, codec.MTU),
		}
		goNet.SessionManager.AddSession(ses)
	} else {
		ses.(*session).conn = conn
	}
	return ses.(*session)
}

// 取原始连接
func (s *session) Socket() interface{} {
	return s.conn
}

func (s *session) Send(msg interface{}) {
	if err := codec.SendPacket(s.conn, msg); err != nil {
		logrus.Errorf("sesssion_%v close error,reason is %v", s.ID(), err)
	}
}

func (s *session) Close() {
	if err := s.conn.Close(); err != nil {
		logrus.Errorf("sesssion_%v close error,reason is %v", s.ID(), err)
	}
	s.store = nil
}

// 接收循环
func (s *session) recvLoop() {
	for {
		n, err := s.conn.Read(s.buf)
		if err != nil {
			logrus.Errorf("session_%v closed,reason is %v", s.ID(), err)
			//recycle session
			goNet.SessionManager.RecycleSession(s)
			break
		}
		controllerIdx, msg, err := codec.ParserPacket(s.buf[:n])
		if err != nil {
			logrus.Warnf("msg parser error,reason is %v", err)
			continue
		}
		controller, err := s.GetController(controllerIdx)
		if err != nil {
			logrus.Warnf("session_%v get controller_%v error, reason is %v", s.ID(), controllerIdx, err)
			continue
		}
		goNet.SubmitMsgToAntsPool(controller, s, msg)
	}
}

func (u *session) Value(v ...interface{}) interface{} {
	if len(v) > 0 {
		u.store = v[0]
	}
	return u.store
}
