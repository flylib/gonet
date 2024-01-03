package udp

import (
	"github.com/flylib/gonet"
	"net"
	"sync"
)

type server struct {
	gonet.PeerCommon
	ln *net.UDPConn
	option
	sync.RWMutex
	remotes map[string]uint64
}

func NewServer(options ...Option) gonet.IServer {
	s := &server{
		option: option{
			mtu: gonet.MTU,
		},
		remotes: map[string]uint64{},
	}
	for _, f := range options {
		f(&s.option)
	}
	return s
}

func (s *server) Listen(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	s.ln, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	s.SetAddr(addr)

	var buf = make([]byte, s.option.mtu)
	for {
		n, remoteAddr, err := s.ln.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		var ses *session
		if sid, exit := s.remotes[remoteAddr.String()]; exit {
			is, _ := s.Context.GetSession(sid)
			ses, _ = is.(*session)
		} else {
			ses = newSession(s.ln, remoteAddr)
			s.remotes[remoteAddr.String()] = ses.ID()
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
