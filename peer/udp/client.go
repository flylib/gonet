package udp

import (
	. "github.com/Quantumoffices/goNet"
	"net"
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
	remoteAddr, err := net.ResolveUDPAddr("udp", c.Addr())
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		panic(err)
	}

	c.session = newSession(conn, nil)
	go c.session.recvLoop()
}

func (c *client) Stop() {
	c.session.Close()
}
