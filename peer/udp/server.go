package udp

import (
	"github.com/sirupsen/logrus"
	"goNet"
	"net"
)

type server struct {
	goNet.PeerIdentify
	session *session
}

func init() {
	identify := goNet.PeerIdentify{}
	identify.SetType("server")
	s := &server{
		PeerIdentify: identify,
	}
	goNet.RegisterPeer(s)
}

func (u *server) Start() {
	localAddr, err := net.ResolveUDPAddr("udp", u.Addr())
	if err != nil {
		logrus.Fatalf("#udp.resolve failed(%s) %v", u.Addr(), err.Error())
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		logrus.Fatalf("#udp.listen failed(%s) %s", u.Addr(), err.Error())
	}
	logrus.Infof("#udp.listen(%s) %s", u.Type(), u.Addr())

	u.session = newSession(conn,localAddr)
	u.session.recvLoop()
}
func (u *server) Stop() {
	u.session.Close()
}
