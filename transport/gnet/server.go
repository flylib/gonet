package gnet

import (
	"context"
	"errors"

	"github.com/flylib/gonet"
	"github.com/panjf2000/gnet/v2"
)

type server struct {
	gnet.EventHandler
	gonet.PeerCommon[*session]
	engine gnet.Engine
	opt    option
}

func NewServer(ctx *gonet.AppContext[*session], options ...Option) gonet.IServer {
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
	ses, ok := is.(*session)
	if !ok {
		return gnet.Close
	}

	var data []byte
	if len(ses.cache) > 0 {
		data = append(ses.cache, buf...)
		ses.cache = ses.cache[:0]
	} else {
		data = buf
	}
	// Loop to handle multiple packets in a single traffic event (TCP粘包).
	for len(data) > 0 {
		msg, unused, err := s.GetCtx().UnPackage(is, data)
		if err != nil {
			if errors.Is(err, gonet.ErrIncompleteHeader) || errors.Is(err, gonet.ErrIncompletePacket) {
				// Cache remaining bytes for the next OnTraffic.
				ses.cache = append(ses.cache[:0], data...)
			} else {
				s.GetCtx().GetLogger().Errorf("gonet gnet: session %d parse error: %v", c.Fd(), err)
			}
			break
		}
		s.GetCtx().PushGlobalMessageQueue(msg)
		if unused <= 0 {
			break
		}
		data = data[len(data)-unused:]
	}
	return gnet.None
}

func (s *server) Listen(addr string) error {
	s.SetAddr(addr)
	return gnet.Run(s, addr, gnet.WithOptions(s.opt.Options))
}

func (s *server) Close() error {
	return s.engine.Stop(context.Background())
}
