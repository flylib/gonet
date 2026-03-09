package gnet

import (
	"context"
	"github.com/flylib/gonet"
	"github.com/panjf2000/gnet/v2"
	"log"
)

type server struct {
	gnet.EventHandler
	gonet.PeerCommon[*session]
	engine gnet.Engine
	opt    option
}

func NewServer(ctx *gonet.Context[*session], options ...Option) gonet.IServer {
	var opt option
	for _, f := range options {
		f(&opt)
	}
	opt.Logger = ctx.GetLogger()

	s := &server{opt: opt}
	s.WithContext(ctx)
	return s
}

// OnBoot fires when the engine is ready for accepting connections.
func (s *server) OnBoot(eng gnet.Engine) (action gnet.Action) {
	s.engine = eng
	return
}

func (s *server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	newSession(s.GetCtx(), c)
	return nil, gnet.None
}

func (s *server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	is, ok := s.GetCtx().GetSession(uint64(c.Fd()))
	if ok {
		s.GetCtx().RecycleSession(is)
	}
	return gnet.None
}

func (s *server) OnTraffic(c gnet.Conn) (action gnet.Action) {
	buf, err := c.Next(-1)
	if err != nil {
		return gnet.Close
	}
	is, ok := s.GetCtx().GetSession(uint64(c.Fd()))
	if !ok {
		return gnet.Close
	}
	message, _, err := s.GetCtx().UnPackage(is, buf)
	if err != nil {
		log.Printf("session_%v msg parser error,reason is %v \n", c.Fd(), err)
		return gnet.None
	}
	s.GetCtx().PushGlobalMessageQueue(message)
	return gnet.None
}

func (s *server) Listen(addr string) error {
	s.SetAddr(addr)
	return gnet.Run(s, addr, gnet.WithOptions(s.opt.Options))
}

func (s *server) Close() error {
	return s.engine.Stop(context.Background())
}
