package tcp

import (
	. "github.com/zjllib/gonet"
	"net"
)

type client struct {
	PeerIdentify
	session *session
}

func (c *client) Start() {
	conn, err := net.Dial("tcp", c.Addr())
	if err != nil {
		panic(err)
	}
	c.session = newSession(conn)
	go c.session.recvLoop()
}

func (c *client) Stop() {
	c.session.Close()
}

func init() {
	identify := PeerIdentify{}
	identify.SetType(PEERTYPE_CLIENT)
	RegisterPeer(&client{PeerIdentify: identify})
}
