package tcp

import (
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
		goNet.Log.Fatalf("tcp(%v) listen failed %v", s.Type(), err.Error())
	}
	s.ln = ln
	goNet.Log.Infof("tcp(%v)listen on %v", s.Type(), ":8087")

	for {
		conn, err := ln.Accept()
		if err != nil {
			goNet.Log.Error("#tcp(%v),accept failed,err:%v", s.Type(), err.Error())
			break
		}
		goNet.Log.Infof("#tcp.accept from %s connected", conn.RemoteAddr())

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
	identify.SetType(goNet.PEERTYPE_SERVER)
	goNet.RegisterPeer(&server{PeerIdentify: identify})
}
