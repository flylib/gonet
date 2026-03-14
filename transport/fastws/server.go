package fastws

import (
	"net/url"

	"github.com/fasthttp/websocket"
	"github.com/flylib/gonet"
	"github.com/valyala/fasthttp"
)

var _ gonet.IServer = new(server)

type server struct {
	gonet.PeerCommon[*Session]
	option
	path string
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
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}
	s.SetAddr(addr)
	s.path = u.Path
	return fasthttp.ListenAndServe(u.Host, s.requestHandler)
}

func (s *server) Close() error {
	return nil
}

func (s *server) requestHandler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != s.path {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	if err := s.option.FastHTTPUpgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		ns := newSession(s.GetCtx(), conn)
		if ns == nil {
			_ = conn.Close()
			return
		}
		ns.ReadLoop()
	}); err != nil {
		return
	}
}
