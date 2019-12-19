package tcp

import (
	"github.com/sirupsen/logrus"
	"goNet"
	"net"
)

type server struct {
	goNet.PeerIdentify
	ln net.Listener
}

func (s *server) Start() {
	ln, err := net.Listen("tcp", s.Addr())
	if err != nil {
		logrus.Fatalf("tcp(%v) listen failed %v", s.Type(), err.Error())
	}
	s.ln = ln
	logrus.Infof("tcp(%v)listen on %v", s.Type(), ":8087")

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Error("#tcp(%v),accept failed,err:%v", s.Type(), err.Error())
			break
		}
		logrus.Infof("#tcp.accept from %s connected", conn.RemoteAddr())

		go s.newConn(conn)
	}
}

//新连接
func (s *server) newConn(conn net.Conn) {
	ses := newSession(conn)
	ses.recvLoop()
}

func (s *server) Stop() {
	s.ln.Close()
}

func init() {
	identify := goNet.PeerIdentify{}
	identify.SetType("server")
	s := &server{
		PeerIdentify: identify,
	}
	goNet.RegisterPeer(s)
}
