package kcp

import (
	"crypto/sha1"
	"golang.org/x/crypto/pbkdf2"

	"reflect"
)

var _ gonet.IServer = new(server)

type server struct {
	gonet.PeerIdentify
	gonet.SessionAbility
	ln *kcp.Listener
}

func NewServer(addr string) *server {
	s := &server{}
	s.SetAddr(addr)
	return s
}

func (s *server) Listen() error {
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	ln, err := kcp.ListenWithOptions(s.Addr(), block, 10, 3)
	if err != nil {
		return err
	}
	s.ln = ln
	for {
		conn, err := s.ln.AcceptKCP()
		if err != nil {
			continue
		}
		go newSession(s.AppContext, conn).recvLoop()
	}
}
func (s *server) Stop() error {
	return s.ln.Close()
}

func (s *server) SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
