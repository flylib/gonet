package ws

import (
	"github.com/gorilla/websocket"
	. "goNet"
	"net/http"
	"time"
)

type client struct {
	PeerIdentify
	session *session
}

func init() {
	identify := PeerIdentify{}
	identify.SetType(PEERTYPE_CLIENT)
	c := &client{
		PeerIdentify: identify,
	}
	RegisterPeer(c)
}

func (c *client) Start() {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}
	conn, _, err := dialer.Dial(c.Addr(), nil)
	if err != nil {
		Log.Panicf("#ws.connect failed(%s) %v", c.Addr(), err.Error())
		return
	}
	Log.Info(conn.RemoteAddr())
	c.session = newSession(conn)
	go c.session.recvLoop()
}

func (c *client) Stop() {
	c.session.conn.SetReadDeadline(time.Now())
}
