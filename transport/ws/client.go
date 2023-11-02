package ws

import (
	"github.com/flylib/gonet"
	"github.com/gorilla/websocket"
	"net/http"
)

type client struct {
	gonet.PeerIdentify
	option
	conn websocket.Conn
}

func NewClient(ctx *gonet.Context, options ...Option) gonet.IClient {
	c := &client{}
	for _, f := range options {
		f(&c.option)
	}
	c.WithContext(ctx)
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
	s := newSession(c.Context, conn)
	go s.ReadLoop()
	return s, nil
}

func (c *client) Close() error {
	return nil
}
