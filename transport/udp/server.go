package udp

import (
	. "github.com/zjllib/gonet/v3"
	"log"
	"net"
	"reflect"
)

var _ IServer = new(server)

type server struct {
	ServerIdentify
	ln *net.UDPConn
}

func NewServer(addr string) *server {
	s := &server{}
	s.SetAddr(addr)
	return s
}

func (s *server) Listen() error {
	addr, err := net.ResolveUDPAddr(string(UDP), s.Addr())
	if err != nil {
		return err
	}
	s.ln, err = net.ListenUDP("udp", addr)
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
			s, _ := s.Context.GetSession(sid)
			ses, _ = s.(*session)
		} else {
			ses = newSession(s.Context, s.ln, remote)
		}
		msg, _, err := s.Context.UnPackage(buf[:n])
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", ses.ID(), err)
			continue
		}
		s.Context.PushGlobalMessageQueue(ses, msg)
	}
	return nil
}
func (s *server) Stop() error {
	return s.ln.Close()
}

func (s *server) SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
