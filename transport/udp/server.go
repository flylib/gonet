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
	conn *net.UDPConn
}

func NewTransport(addr string) *server {
	s := &server{}
	s.SetAddr(addr)
	return s
}

func (s *server) Listen() error {
	localAddr, err := net.ResolveUDPAddr(string(UDP), s.Addr())
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return err
	}
	s.conn = conn
	for {
		var buf []byte
		n, remote, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("#udp.read failed(%v) %v \n", s.conn.RemoteAddr(), err.Error())
			continue
		}
		var ses *session
		if sid, exit := remotes[remote.String()]; exit {
			s, _ := GetSession(sid)
			ses, _ = s.(*session)
		} else {
			ses = newSession(s.conn, remote)
		}
		msg, _, err := ParserPacket(buf[:n])
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", ses.ID(), err)
			continue
		}
		HandingMessage(ses, msg)
	}
	return nil
}
func (s *server) Stop() error {
	return s.conn.Close()
}

func (s *server) SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
