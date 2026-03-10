package udp

import (
	"github.com/flylib/gonet"
	"net"
	"sync"
)

type server struct {
	gonet.PeerCommon[*session]
	ln *net.UDPConn
	option
	sync.RWMutex
	remotes map[string]uint64
}

func NewServer(ctx *gonet.AppContext[*session], options ...Option) gonet.IServer {
	s := &server{
		option: option{
			mtu: gonet.MTU,
		},
		remotes: map[string]uint64{},
	}
	for _, f := range options {
		f(&s.option)
	}
	s.WithContext(ctx)
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
			s.GetCtx().GetLogger().Errorf("#udp.read failed(%v) %v \n", s.ln.RemoteAddr(), err.Error())
			return err
		}

		var ses *session
		if sid, exist := s.remotes[remoteAddr.String()]; exist {
			is, _ := s.GetCtx().GetSession(sid)
			ses, _ = is.(*session)
		} else {
			ses = newSession(s.GetCtx(), s.ln, remoteAddr)
			if ses == nil {
				continue
			}
			s.remotes[remoteAddr.String()] = ses.ID()
		}
		msg, _, err := s.GetCtx().UnPackage(ses, buf[:n])
		if err != nil {
			s.GetCtx().GetLogger().Errorf("session_%v msg parser error,reason is %v \n", ses.ID(), err)
			continue
		}
		s.GetCtx().PushGlobalMessageQueue(msg)
	}
}

func (s *server) Close() error {
	return s.ln.Close()
}
