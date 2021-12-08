package tcp

import (
	. "github.com/zjllib/gonet/v3"
	"net"
)

type server struct {
	ServerIdentify
	ln net.Listener
}

func init() {
	RegisterServer(&server{}, session{})
}

func (s *server) Start() error {
	ln, err := net.Listen(string(TCP), s.Addr())
	if err != nil {
		return err
	}
	s.ln = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.newConn(conn)
	}
}
func (s *server) Stop() error {
	return s.ln.Close()
}

//新连接
func (s *server) newConn(conn net.Conn) {
	ses := newSession(conn)
	ses.recvLoop()
}
