package gnet

import (
	"context"
	"github.com/flylib/gonet"
	"github.com/panjf2000/gnet/v2"
	"log"
	"reflect"
)

type server struct {
	gnet.EventHandler
	gonet.PeerIdentify
	engine gnet.Engine
	opt    option
}

func NewServer(ctx *gonet.Context, options ...Option) gonet.IServer {
	var opt option
	for _, f := range options {
		f(&opt)
	}
	opt.Logger = ctx.ILogger

	s := &server{opt: opt}
	s.WithContext(ctx)
	ctx.InitSessionMgr(reflect.TypeOf(session{}))
	return s
}

// OnBoot fires when the engine is ready for accepting connections.
// The parameter engine has information and various utilities.
func (s *server) OnBoot(eng gnet.Engine) (action gnet.Action) {
	s.engine = eng
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
	is, ok := s.Context.GetSession(uint64(c.Fd()))
	if !ok {
		return gnet.Close
	}
	message, _, err := s.Context.UnPackage(is, buf)
	if err != nil {
		log.Printf("session_%v msg parser error,reason is %v \n", c.Fd(), err)
		return gnet.None
	}
	s.Context.PushGlobalMessageQueue(message)
	return gnet.None
}

func (s *server) Listen(addr string) error {
	s.SetAddr(addr)
	return gnet.Run(s, addr, gnet.WithOptions(s.opt.Options))
}

func (s *server) Close() error {
	return s.engine.Stop(context.Background())
}
