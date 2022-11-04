package tcp

import (
	"github.com/zjllib/gonet/v3/transport"
	"net"
	"reflect"
)

var _ transport.Transport = new(tcp)

type tcp struct {
	transport.TransportIdentify
	ln net.Listener
}

func NewTransport(addr string) *tcp {
	s := &tcp{}
	s.SetAddr(addr)
	return s
}

func (s *tcp) Listen() error {
	ln, err := net.Listen(string(transport.TCP), s.Addr())
	if err != nil {
		return err
	}
	s.ln = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go newSession(conn).recvLoop()
	}
}
func (s *tcp) Stop() error {
	return s.ln.Close()
}

func (s *tcp) SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
