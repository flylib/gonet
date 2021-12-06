package udp

import (
	"github.com/astaxie/beego/logs"
	"github.com/zjllib/gonet/v3/transport"
	"net"
)

//addr:sessionID
var remotes = map[string]uint32{}

// Socket会话
type session struct {
	SessionIdentify
	remote *net.UDPAddr
	conn   *net.UDPConn
	data   interface{}
	buf    []byte
}

//新会话
func newSession(conn *net.UDPConn, remote *net.UDPAddr) *session {
	ses := sessions.GetIdleSession()
	if ses == nil {
		ses = &session{
			conn:   conn,
			remote: remote,
			buf:    make([]byte, transport.MTU),
		}
		sessions.AddSession(ses)
	} else {
		ses.(*session).conn = conn
	}
	remotes[remote.String()] = ses.ID()
	return ses.(*session)
}

// 取原始连接
func (s *session) Socket() interface{} {
	return s.conn
}

// 发送封包
func (s *session) Send(msg interface{}) {
	var err error
	if s.remote == nil {
		logs.Info("client send msg ")
		err = transport.SendPacket(s.conn, msg)
	} else {
		logs.Info("server send msg ")
		err = transport.SendUdpPacket(s.conn, msg, s.remote)
	}
	if err != nil {
		logs.Error("sesssion_%v close error,reason is %v", s.ID(), err)
	}
}

func (s *session) Close() {
	if err := s.conn.Close(); err != nil {
		logs.Error("sesssion_%v close error,reason is %v", s.ID(), err)
	}
	s.data = nil
}

// 接收循环
func (s *session) recvLoop() {
	for {
		n, remote, err := s.conn.ReadFromUDP(s.buf)
		logs.Info("recv=", remote.String())
		if err != nil {
			logs.Errorf("#udp.accept failed(%v) %v", s.conn.RemoteAddr(), err.Error())
		}
		var ses session
		if sid, exit := remotes[remote.String()]; exit {
			ses = SessionManager.GetSessionById(sid)
		} else {
			ses = newSession(s.conn, remote)
		}
		//msg, err := codec.ParserPacket(s.buf[:n])
		//if err != nil {
		//	logs.Warnf("message decode error=%s", err)
		//	continue
		//}
		//SubmitMsgToAntsPool(msg, ses)
	}
}

func (u *session) Value(v ...interface{}) interface{} {
	if len(v) > 0 {
		u.data = v[0]
	}
	return u.data
}
