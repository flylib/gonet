package tcp

import (
	"github.com/flylib/gonet"
	"net"
)

type client struct {
	gonet.PeerCommon[*Session]
	option
}

func NewClient(ctx *gonet.Context[*Session], options ...Option) gonet.IClient {
	c := &client{}
	for _, f := range options {
		f(&c.option)
	}
	c.WithContext(ctx)
	return c
}

func (c *client) Dial(addr string) (gonet.ISession, error) {
	var conn net.Conn
	var err error
	if c.option.HandshakeTimeout > 0 {
		conn, err = net.DialTimeout("tcp", addr, c.option.HandshakeTimeout)
	} else {
		conn, err = net.Dial("tcp", addr)
	}
	if err != nil {
		return nil, err
	}
	c.SetAddr(addr)
	s := newSession(c.GetCtx(), conn)
	if s == nil {
		_ = conn.Close()
		return nil, nil
	}
	go s.recvLoop()
	return s, nil
}

func (c *client) Close() error {
	return nil
}
