package udp

import (
	"github.com/flylib/gonet"
	"net"
)

var _ gonet.IClient = new(client)

type client struct {
	gonet.PeerCommon[*session]
	conn *net.UDPConn
	option
}

func NewClient(ctx *gonet.Context[*session], options ...Option) gonet.IClient {
	c := &client{
		option: option{
			mtu: gonet.MTU,
		},
	}
	for _, f := range options {
		f(&c.option)
	}
	c.WithContext(ctx)
	return c
}

func (c *client) Dial(addr string) (gonet.ISession, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	c.SetAddr(addr)
	c.conn = conn

	s := newSession(c.GetCtx(), conn, udpAddr)
	if s == nil {
		_ = conn.Close()
		return nil, nil
	}
	s.remoteConn = conn
	go s.recvLoop()
	return s, nil
}

func (c *client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
