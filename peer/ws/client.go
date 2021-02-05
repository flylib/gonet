package ws

import (
	. "github.com/Quantumoffices/goNet"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
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
	logs.Info(c.Addr())
	conn, _, err := dialer.Dial(c.Addr(), nil)
	if err != nil {
		panic(err)
	}
	c.session = newSession(conn)
	go c.session.recvLoop()
	go c.session.sendLoop()
}

func (c *client) Stop() {
	c.session.Conn.SetReadDeadline(time.Now())
}
