package udp

import (
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/transport"
	"log"
	"net"
)

type server struct {
	ServerIdentify
	conn *net.UDPConn
}

func init() {
	RegisterServer(&server{}, session{})
}

func (u *server) Start() error {
	localAddr, err := net.ResolveUDPAddr("udp", u.Addr())
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return err
	}
	u.conn = conn
	//u.session = newSession(conn, localAddr)
	for {
		var buf []byte
		n, remote, err := u.conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("#udp.read failed(%v) %v \n", u.conn.RemoteAddr(), err.Error())
			continue
		}
		var ses *session
		if sid, exit := remotes[remote.String()]; exit {
			s, _ := GetSession(sid)
			ses, _ = s.(*session)
		} else {
			ses = newSession(u.conn, remote)
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
func (u *server) Stop() error {
	return u.conn.Close()
}
