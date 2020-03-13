package tcp

import (
	"goNet"
	"net"
)

type client struct {
	goNet.PeerIdentify
	session *session
}

func (c *client) Start() {
	conn, err := net.Dial("tcp", c.Addr())
	if err != nil {
		goNet.Log.Fatalf("#tcp(%v) connect failed %v", c.Type(), err.Error())
		return
	}
	goNet.Log.Infof("#tcp(%v) connect(%v) success", c.Type(), conn.RemoteAddr())
	c.session = newSession(conn)
	go c.session.recvLoop()
}

func (c *client) Stop() {
	c.session.Close()
}

func init() {
	identify := goNet.PeerIdentify{}
	identify.SetType(goNet.PEERTYPE_CLIENT)
	goNet.RegisterPeer(&client{PeerIdentify: identify})
}
