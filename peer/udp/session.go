package udp

import (
	"github.com/sirupsen/logrus"
	"goNet"
	"goNet/codec"
	"net"
)

//addr:sessionID
var remotes = map[string]uint32{}

// Socket会话
type session struct {
	goNet.SessionIdentify
	remote *net.UDPAddr
	conn   *net.UDPConn
	data   interface{}
	buf    []byte
}

//新会话
func newSession(conn *net.UDPConn, remote *net.UDPAddr) *session {
	ses := goNet.SessionManager.GetIdleSession()
	if ses == nil {
		ses = &session{
			conn:   conn,
			remote: remote,
			buf:    make([]byte, codec.MTU),
		}
		goNet.SessionManager.AddSession(ses)
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
		logrus.Info("client send msg ")
		err = codec.SendPacket(s.conn, msg)
	} else {
		logrus.Info("server send msg ")
		err = codec.SendUdpPacket(s.conn, msg, s.remote)
	}
	if err != nil {
		logrus.Errorf("sesssion_%v close error,reason is %v", s.ID(), err)
	}
}

func (s *session) Close() {
	if err := s.conn.Close(); err != nil {
		logrus.Errorf("sesssion_%v close error,reason is %v", s.ID(), err)
	}
	s.data = nil
}

// 接收循环
func (s *session) recvLoop() {
	for {
		n, remote, err := s.conn.ReadFromUDP(s.buf)
		logrus.Info("recv=", remote.String())
		if err != nil {
			logrus.Errorf("#udp.accept failed(%v) %v", s.conn.RemoteAddr(), err.Error())
		}
		var ses goNet.Session
		if sid, exit := remotes[remote.String()]; exit {
			ses = goNet.SessionManager.GetSessionById(sid)
		} else {
			ses = newSession(s.conn, remote)
		}
		msg, err := codec.ParserPacket(s.buf[:n])
		if err != nil {
			logrus.Warnf("message decode error=%s", err)
			continue
		}
		goNet.HandleMessage(msg, ses)
	}
}

func (u *session) Value(v ...interface{}) interface{} {
	if len(v) > 0 {
		u.data = v[0]
	}
	return u.data
}
