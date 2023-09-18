package ws

import (
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"net/http"
	"time"
)

type client struct {
	PeerIdentify
}

func (c *client) Dial() (ISession, error) {
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
