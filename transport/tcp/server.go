package tcp

import (
	"github.com/flylib/gonet"
	"net"
)

type server struct {
	gonet.PeerCommon[*Session]
	ln net.Listener
}

func NewServer(ctx *gonet.Context[*Session]) gonet.IServer {
	s := &server{}
	s.WithContext(ctx)
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
		session := newSession(s.GetCtx(), conn)
		if session == nil {
			// max session count reached
			_ = conn.Close()
			continue
		}
		go session.recvLoop()
	}
}

func (s *server) Close() error {
	return s.ln.Close()
}
