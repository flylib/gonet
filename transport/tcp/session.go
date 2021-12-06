package tcp

import (
	"github.com/astaxie/beego/logs"
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/codec"
	"github.com/zjllib/gonet/v3/transport"
	"net"
)

// Socket会话
type session struct {
	//核心会话标志
	SessionIdentify
	//存储功能
	SessionStore
	//累计收消息总数
	recvCount uint64
	//raw conn
	conn net.Conn
	buf  []byte
	//缓存数据，用于解决粘包问题
	//cache []byte
}

//新会话
func newSession(conn net.Conn) *session {
	ses := CreateSession()
	ses.(*session).conn = conn
	return ses.(*session)
}

// 取原始连接
func (s *session) Socket() interface{} {
	return s.conn
}

func (s *session) Send(msg interface{}) error {
	return transport.SendPacket(s.conn, msg)
}

func (s *session) Close() error {
	if err := s.conn.Close(); err != nil {
		return err
	}
	s.SessionStore = SessionStore{}
	return nil
}

// 接收循环
func (s *session) recvLoop() {
	for {
		n, err := s.conn.Read(s.buf)
		if err != nil {
			logs.Error("session_%v closed,reason is %v", s.ID(), err)
			//recycle session
			RecycleSession(s)
			break
		}
		actorIDx, msg, err := codec.ParserPacket(s.buf[:n])
		if err != nil {
			logs.Warn("msg parser error,reason is %v", err)
			continue
		}
		controller, err := s.GetController(actorIDx)
		if err != nil {
			logs.Warn("session_%v get controller_%v error, reason is %v", s.ID(), actorIDx, err)
			continue
		}
		HandleEvent(controller, s, msg)
	}
}
