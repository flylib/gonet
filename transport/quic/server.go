package quic

import (
	"context"
	"github.com/flylib/gonet"
	"github.com/quic-go/quic-go"
)

// 接收端
type server struct {
	gonet.PeerCommon
	ln *quic.Listener
	option
}

func NewServer(ctx *gonet.Context, options ...Option) gonet.IServer {
	s := &server{}
	for _, f := range options {
		f(&s.option)
	}
	s.WithContext(ctx)
	return s
}

func (s *server) Listen(addr string) (err error) {
	s.ln, err = quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	s.SetAddr(addr)

	for {
		conn, err := s.ln.Accept(context.Background())
		if err != nil {
			continue
		}
		ns := newSession(s.Context, conn)
		ns.mod = s.modulo
		go ns.acceptStream()
	}
}

func (s *server) Close() error {
	return s.ln.Close()
}
