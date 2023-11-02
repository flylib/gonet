package tcp

import (
	"github.com/flylib/gonet"
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
	if c.option.HandshakeTimeout > 0 {
		conn, err := net.DialTimeout("tcp", addr, c.option.HandshakeTimeout)
		if err != nil {
			return nil, err
		}
		c.conn = conn
	} else {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		c.conn = conn
	}
	s := newSession(c.Context, c.conn)
	go s.recvLoop()
	return s, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
