package udp

import (
	. "github.com/zjllib/gonet/v3"
	"net"
	"reflect"
)

type client struct {
	PeerIdentify
}

func NewClient(ctx *Context) IClient {
	c := &client{}
	c.WithContext(ctx)
	ctx.InitSessionMgr(reflect.TypeOf(session{}))
	return c
}

func (c *client) Dial(addr string) (ISession, error) {
	c.SetAddr(addr)

	conn, err := net.Dial("udp", "127.0.0.1:9001")
	if err != nil {
		return nil, err
	}

	udpAddr, err := net.ResolveUDPAddr(string(UDP), c.Addr())
	if err != nil {
		return nil, err
	}
	s := newSession(c.Context, conn, udpAddr)
	go s.()
	return s, nil
}
