package ws

import (
	"github.com/gorilla/websocket"
	"goNet"
	"net/http"
	"time"
)

type client struct {
	goNet.PeerIdentify
	session *session
}

func init() {
	identify := goNet.PeerIdentify{}
	identify.SetType(goNet.PEERTYPE_CLIENT)
	c := &client{
		PeerIdentify: identify,
	}
	goNet.RegisterPeer(c)
}

func (c *client) Start() {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}
	conn, _, err := dialer.Dial(c.Addr(), nil)
	if err != nil {
		goNet.Log.Errorf("#ws.connect failed(%s) %v", c.Addr(), err.Error())
		return
	}
	goNet.Log.Info(conn.RemoteAddr())
	c.session = newSession(conn)
	go c.session.recvLoop()
}

func (c *client) Stop() {
	c.session.conn.SetReadDeadline(time.Now())
}
