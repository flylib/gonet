package udp

import (
	"github.com/flylib/gonet"
	"net"
)

var _ gonet.IClient = new(client)

type client struct {
	gonet.PeerCommon
	conn *net.UDPConn
	option
}

func NewClient(options ...Option) gonet.IClient {
	c := &client{
		option: option{
			mtu: gonet.MTU,
		},
	}
	for _, f := range options {
		f(&c.option)
	}
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

	s := newSession(conn, udpAddr)
	s.remoteConn = conn
	go s.readLoop()
	return s, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
