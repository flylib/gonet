package udp

import (
	. "github.com/Quantumoffices/goNet"
	"net"
)

type server struct {
	PeerIdentify
	session *session
}

func init() {
	identify := PeerIdentify{}
	identify.SetType(PEERTYPE_SERVER)
	s := &server{
		PeerIdentify: identify,
	}
	RegisterPeer(s)
}

func (u *server) Start() {
	localAddr, err := net.ResolveUDPAddr("udp", u.Addr())
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		panic(err)
	}
	u.session = newSession(conn, localAddr)
	u.session.recvLoop()
}
func (u *server) Stop() {
	u.session.Close()
}
