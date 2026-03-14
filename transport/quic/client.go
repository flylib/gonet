package quic

import (
	"context"
	"github.com/flylib/gonet"
	"github.com/quic-go/quic-go"
)

type client struct {
	gonet.PeerCommon[*session]
	option
	conn quic.Connection
}

func NewClient(ctx *gonet.AppContext[*session], options ...Option) gonet.IClient {
	c := &client{option: option{modulo: defaultChannelModulo}}
	for _, f := range options {
		f(&c.option)
	}
	if c.option.modulo == 0 {
		c.option.modulo = defaultChannelModulo
	}
	c.WithContext(ctx)
	return c
}

func (c *client) Dial(addr string) (gonet.ISession, error) {
	connection, err := quic.DialAddr(context.Background(), addr, generateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}
	c.conn = connection
	s := newSession(c.GetCtx(), connection)
	if s == nil {
		_ = connection.CloseWithError(0, "max sessions reached")
		return nil, nil
	}
	s.mod = c.modulo
	go s.acceptStream()
	return s, nil
}

func (c *client) Close() error {
	return c.conn.CloseWithError(0, "EOF")
}
