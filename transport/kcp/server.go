package kcp

import (
	"crypto/sha1"
	"github.com/flylib/gonet"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

type server struct {
	gonet.PeerCommon[*Session]
	ln *kcp.Listener
	option
}

func NewServer(ctx *gonet.AppContext[*Session], options ...Option) gonet.IServer {
	s := &server{}
	for _, f := range options {
		f(&s.option)
	}
	s.WithContext(ctx)
	return s
}

func (s *server) Listen(addr string) error {
	key := pbkdf2.Key([]byte(s.PBKDF2Password), []byte(s.PBKDF2Salt), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	ln, err := kcp.ListenWithOptions(addr, block, 10, 3)
	if err != nil {
		return err
	}
	s.ln = ln
	s.SetAddr(addr)

	for {
		conn, err := s.ln.AcceptKCP()
		if err != nil {
			continue
		}
		session := newSession(s.GetCtx(), conn)
		if session == nil {
			_ = conn.Close()
			continue
		}
		go session.recvLoop()
	}
}

func (s *server) Close() error {
	return s.ln.Close()
}
