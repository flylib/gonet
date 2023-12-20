package tcp

import (
	"github.com/flylib/gonet"
	"net"
)

type server struct {
	gonet.PeerCommon

	ln net.Listener
}

func NewServer(ctx *gonet.Context) gonet.IServer {
	s := &server{}
	s.WithContext(ctx)
	return s
}

func (s *server) Listen(addr string) error {
	ln, err := net.Listen(string(gonet.TCP), addr)
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
		go newSession(s.Context, conn).recvLoop()
	}
}
func (s *server) Close() error {
	return s.ln.Close()
}
