package tcp

import (
	"github.com/zjllib/gonet/v3"
	"net"
	"reflect"
)

var _ gonet.IServer = new(server)

type server struct {
	gonet.ServerIdentify
	gonet.SessionStore
	ln net.Listener
}

func NewServer(addr string) *server {
	s := &server{}
	s.SetAddr(addr)
	return s
}

func (s *server) Listen() error {
	ln, err := net.Listen(string(gonet.TCP), s.Addr())
	if err != nil {
		return err
	}
	s.ln = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go newSession(s.Context, conn).recvLoop()
	}
}
func (s *server) Stop() error {
	return s.ln.Close()
}

func (s *server) SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
