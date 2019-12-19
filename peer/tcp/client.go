package tcp

import (
	"github.com/sirupsen/logrus"
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
		logrus.Fatalf("#tcp(%v) connect failed %v", c.Type(), err.Error())
		return
	}
	logrus.Infof("#tcp(%v) connect(%v) success", c.Type(), conn.RemoteAddr())
	c.session = newSession(conn)
	go c.session.recvLoop()
}

func (c *client) Stop() {
	c.session.Close()
}

func init() {
	identify := goNet.PeerIdentify{}
	identify.SetType("client")
	c := &client{
		PeerIdentify: identify,
	}
	goNet.RegisterPeer(c)
}
