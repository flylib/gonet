package tcp

import (
	"github.com/astaxie/beego/logs"
	. "github.com/zjllib/goNet"
	"net"
)

type server struct {
	PeerIdentify
	ln net.Listener
}

func (s *server) Start() {
	ln, err := net.Listen("tcp", s.Addr())
	if err != nil {
		panic(err)
	}
	s.ln = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Error("#tcp(%v),accept failed,err:%v", s.Type(), err.Error())
			break
		}
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
	identify := PeerIdentify{}
	identify.SetType(PeertypeServer)
	RegisterPeer(&server{PeerIdentify: identify})
}
