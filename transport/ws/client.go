package ws

import (
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"net/http"
	"reflect"
	"time"
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
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}
	conn, _, err := dialer.Dial(c.Addr(), nil)
	if err != nil {
		return nil, err
	}
	s := newSession(c.Context, conn)
	go s.readLoop()
	return s, nil
}
