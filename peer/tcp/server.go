package tcp

import (
	"goNet"
	. "goNet/log"
	"net"
)

type server struct {
	goNet.PeerIdentify
	ln net.Listener
}

func (s *server) Start() {
	ln, err := net.Listen("tcp", s.Addr())
	if err != nil {
		Log.Fatalf("tcp(%v) listen failed %v", s.Type(), err.Error())
	}
	s.ln = ln
	Log.Infof("tcp(%v)listen on %v", s.Type(), ":8087")

	for {
		conn, err := ln.Accept()
		if err != nil {
			Log.Error("#tcp(%v),accept failed,err:%v", s.Type(), err.Error())
			break
		}
		Log.Infof("#tcp.accept from %s connected", conn.RemoteAddr())

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
	identify.SetType(goNet.PEER_SERVER)
	goNet.RegisterPeer(&server{PeerIdentify: identify})
}
