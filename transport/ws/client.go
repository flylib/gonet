package ws

import (
	"github.com/gorilla/websocket"
	"github.com/zjllib/gonet/v3"
	"net/http"
	"reflect"
)

type client struct {
	gonet.PeerIdentify
	option
}

func NewClient(ctx *gonet.AppContext, options ...Option) gonet.IClient {
	c := &client{}
	for _, f := range options {
		f(&c.option)
	}
	c.WithContext(ctx)
	ctx.InitSessionMgr(reflect.TypeOf(session{}))
	return c
}

func (c *client) Dial(addr string) (gonet.ISession, error) {
	c.SetAddr(addr)
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: c.option.HandshakeTimeout,
	}
	conn, _, err := dialer.Dial(c.Addr(), nil)
	if err != nil {
		return nil, err
	}
	s := newSession(c.AppContext, conn)
	go s.ReadLoop()
	go s.SendLoop(s.write)
	return s, nil
}
