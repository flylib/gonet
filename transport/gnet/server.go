package gnet

import (
	"context"
	"github.com/panjf2000/gnet/v2"
	"github.com/zjllib/gonet/v3"
	"log"
	"reflect"
)

var _ gonet.IServer = new(server)

type server struct {
	gnet.EventHandler
	gonet.PeerIdentify
	ln gnet.Engine
}

func NewServer(addr string) *server {
	s := &server{}
	s.SetAddr(addr)
	return s
}

// OnBoot fires when the engine is ready for accepting connections.
// The parameter engine has information and various utilities.
func (s *server) OnBoot(eng gnet.Engine) (action gnet.Action) {
	s.ln = eng
	return
}
func (s *server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	newSession(s.Context, c)
	return nil, gnet.None
}
func (s *server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	is, ok := s.Context.GetSession(uint64(c.Fd()))
	if ok {
		s.Context.RecycleSession(is, err)
	}
	return gnet.None
}

func (s *server) OnTraffic(c gnet.Conn) (action gnet.Action) {
	buf, err := c.Next(-1)
	if err != nil {
		return gnet.Close
	}
	message, _, err := s.Context.UnPackage(buf)
	if err != nil {
		log.Printf("session_%v msg parser error,reason is %v \n", c.Fd(), err)
		return gnet.None
	}
	is, ok := s.Context.GetSession(uint64(c.Fd()))
	if ok {
		s.Context.PushGlobalMessageQueue(is, message)
	}
	return gnet.None
}

func (s *server) Listen() error {
	return gnet.Run(s, s.Addr(), gnet.WithMulticore(true))
}
func (s *server) Stop() error {
	return s.ln.Stop(context.Background())
}

func (s *server) SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
