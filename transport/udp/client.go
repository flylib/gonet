package udp

import (
	"github.com/flylib/gonet"
	"net"
)

var _ gonet.IClient = new(client)

type client struct {
	gonet.PeerIdentify
	conn *net.UDPConn
	option
}

func NewClient(ctx *gonet.Context, options ...Option) gonet.IClient {
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
	udpAddr, err := net.ResolveUDPAddr(string(gonet.UDP), addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP(string(gonet.UDP), nil, udpAddr)
	if err != nil {
		return nil, err
	}
	c.SetAddr(addr)

	s := newSession(c.Context, conn, udpAddr)
	s.remoteConn = conn
	go s.recvLoop()
	return s, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
