package udp

import (
	. "goNet"
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
		Log.Fatalf("#resolve udp address failed(%s) %v", c.Addr(), err.Error())
	}

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		Log.Fatalf("#udp.connect failed(%s) %v", c.Addr(), err.Error())
	}

	c.session = newSession(conn, nil)
	go c.session.recvLoop()
}

func (c *client) Stop() {
	c.session.Close()
}
