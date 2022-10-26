package udp

import (
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/transport"
	"log"
	"net"
	"reflect"
)

type udp struct {
	transport.TransportIdentify
	conn *net.UDPConn
}

func NewTransport(addr string) *udp {
	s := &udp{}
	s.SetAddr(addr)
	return s
}

func (s *udp) Listen() error {
	localAddr, err := net.ResolveUDPAddr(string(transport.UDP), s.Addr())
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
		msg, _, err := transport.ParserPacket(buf[:n])
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", ses.ID(), err)
			continue
		}
		msg.Session = ses
		CacheMsg(msg)
	}
	return nil
}
func (s *udp) Stop() error {
	return s.conn.Close()
}

func (s *udp) SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
