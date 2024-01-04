package kcp

import (
	"crypto/sha1"
	"github.com/flylib/gonet"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"reflect"
)

type server struct {
	gonet.PeerCommon

	ln *kcp.Listener
	option
}

func NewServer(options ...Option) gonet.IServer {
	s := &server{}
	for _, f := range options {
		f(&s.option)
	}
	return s
}

func (s *server) Listen(url string) error {
	key := pbkdf2.Key([]byte(s.PBKDF2Password), []byte(s.PBKDF2Salt), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	ln, err := kcp.ListenWithOptions(url, block, 10, 3)
	if err != nil {
		return err
	}
	s.ln = ln

	s.SetAddr(url)

	for {
		conn, err := s.ln.AcceptKCP()
		if err != nil {
			continue
		}
		go newSession(conn).recvLoop()
	}
}

func (s *server) Close() error {
	return s.ln.Close()
}

func (s *server) SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
