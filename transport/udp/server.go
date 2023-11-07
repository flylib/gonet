package udp

import (
	"github.com/flylib/gonet"
	"net"
)

type server struct {
	gonet.PeerIdentify
	ln *net.UDPConn
	option
}

func NewServer(ctx *gonet.Context, options ...Option) gonet.IServer {
	s := &server{
		option: option{
			mtu: gonet.MTU,
		},
	}
	for _, f := range options {
		f(&s.option)
	}
	s.WithContext(ctx)
	return s
}

func (s *server) Listen(addr string) error {
	udpAddr, err := net.ResolveUDPAddr(string(gonet.UDP), addr)
	if err != nil {
		return err
	}
	s.ln, err = net.ListenUDP(string(gonet.UDP), udpAddr)
	if err != nil {
		return err
	}
	s.SetAddr(addr)

	var buf = make([]byte, s.option.mtu)
	for {
		n, remoteAddr, err := s.ln.ReadFromUDP(buf)
		if err != nil {
			s.ILogger.Errorf("#udp.read failed(%v) %v \n", s.ln.RemoteAddr(), err.Error())
			return err
		}

		var ses *session
		if sid, exit := remotes[remoteAddr.String()]; exit {
			is, _ := s.Context.GetSession(sid)
			ses, _ = is.(*session)
		} else {
			ses = newSession(s.Context, s.ln, remoteAddr)
		}

		msg, _, err := s.Context.UnPackage(ses, buf[:n])
		if err != nil {
			s.ILogger.Errorf("session_%v msg parser error,reason is %v \n", ses.ID(), err)
			continue
		}
		s.Context.PushGlobalMessageQueue(msg)
	}
	return nil
}

func (s *server) Close() error {
	return s.ln.Close()
}
