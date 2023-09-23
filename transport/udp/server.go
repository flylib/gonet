package udp

import (
	. "github.com/zjllib/gonet/v3"
	"log"
	"net"
	"reflect"
)

var _ IServer = new(server)

type server struct {
	PeerIdentify
	ln *net.UDPConn
}

func NewServer(ctx *AppContext) IServer {
	s := &server{}
	s.WithContext(ctx)
	ctx.InitSessionMgr(reflect.TypeOf(session{}))
	return s
}
func (s *server) Listen(addr string) error {
	s.SetAddr(addr)
	udpAddr, err := net.ResolveUDPAddr(string(UDP), s.Addr())
	if err != nil {
		return err
	}
	s.ln, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	for {
		var buf = make([]byte, MTU)
		n, remote, err := s.ln.ReadFromUDP(buf)
		if err != nil {
			log.Printf("#udp.read failed(%v) %v \n", s.ln.RemoteAddr(), err.Error())
			return err
		}
		var ses *session
		if sid, exit := remotes[remote.String()]; exit {
			s, _ := s.AppContext.GetSession(sid)
			ses, _ = s.(*session)
		} else {
			ses = newSession(s.AppContext, s.ln, remote)
		}
		msg, _, err := s.AppContext.UnPackage(buf[:n])
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", ses.ID(), err)
			continue
		}
		s.AppContext.PushGlobalMessageQueue(ses, msg)
	}
	return nil
}
func (s *server) Stop() error {
	return s.ln.Close()
}
