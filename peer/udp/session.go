package udp

import (
	. "goNet"
	"goNet/codec"
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
	ses := SessionManager.GetIdleSession()
	if ses == nil {
		ses = &session{
			conn:   conn,
			remote: remote,
			buf:    make([]byte, codec.MTU),
		}
		SessionManager.AddSession(ses)
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
		Log.Info("client send msg ")
		err = codec.SendPacket(s.conn, msg)
	} else {
		Log.Info("server send msg ")
		err = codec.SendUdpPacket(s.conn, msg, s.remote)
	}
	if err != nil {
		Log.Errorf("sesssion_%v close error,reason is %v", s.ID(), err)
	}
}

func (s *session) Close() {
	if err := s.conn.Close(); err != nil {
		Log.Errorf("sesssion_%v close error,reason is %v", s.ID(), err)
	}
	s.data = nil
}

// 接收循环
func (s *session) recvLoop() {
	//for {
	//	n, remote, err := s.conn.ReadFromUDP(s.buf)
	//	Log.Info("recv=", remote.String())
	//	if err != nil {
	//		Log.Errorf("#udp.accept failed(%v) %v", s.conn.RemoteAddr(), err.Error())
	//	}
	//	var ses Session
	//	if sid, exit := remotes[remote.String()]; exit {
	//		ses = SessionManager.GetSessionById(sid)
	//	} else {
	//		ses = newSession(s.conn, remote)
	//	}
	//	//msg, err := codec.ParserPacket(s.buf[:n])
	//	//if err != nil {
	//	//	Log.Warnf("message decode error=%s", err)
	//	//	continue
	//	//}
	//	//SubmitMsgToAntsPool(msg, ses)
	//}
}

func (u *session) Value(v ...interface{}) interface{} {
	if len(v) > 0 {
		u.data = v[0]
	}
	return u.data
}
