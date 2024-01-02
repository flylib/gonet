package tcp

import (
	"github.com/flylib/gonet"
	"net"
)

type server struct {
	gonet.PeerCommon

	ln net.Listener

	option
}

func NewServer(options ...Option) gonet.IServer {
	s := &server{}
	for _, f := range options {
		f(&s.option)
	}
	return s
}

func (s *server) Listen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.ln = ln
	s.SetAddr(addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go newSession(conn).recvLoop()
	}
}
func (s *server) Close() error {
	return s.ln.Close()
}
