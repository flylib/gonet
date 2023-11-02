package quic

import (
	"github.com/flylib/gonet"
	"github.com/quic-go/quic-go"
	"net"
)

type client struct {
	gonet.PeerIdentify
	option
	conn net.Conn
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
	quic.DialAddr()
}

func (c *client) Close() error {
	return c.conn.Close()
}
