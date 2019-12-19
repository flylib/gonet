package udp

import (
	"goNet"
	"log"
	"net"
)

type client struct {
	goNet.PeerIdentify
	session *session
}

func init() {
	identify := goNet.PeerIdentify{}
	identify.SetType("client")
	c := &client{
		PeerIdentify: identify,
	}
	goNet.RegisterPeer(c)
}

func (c *client) Start() {
	remoteAddr, err := net.ResolveUDPAddr("udp", c.Addr())
	if err != nil {
		log.Fatalf("#resolve udp address failed(%s) %v", c.Addr(), err.Error())
	}

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Fatalf("#udp.connect failed(%s) %v", c.Addr(), err.Error())
	}

	c.session = newSession(conn,nil)
	go c.session.recvLoop()
}

func (c *client) Stop() {
	c.session.Close()
}
