package gonet

type IPeerIdentify interface {
	Addr() string
	SetAddr(addr string)
	WithContext(c *AppContext)
}
type PeerIdentify struct {
	*AppContext
	uuid string
	addr string
}

func (s *PeerIdentify) Addr() string {
	return s.addr
}

func (s *PeerIdentify) SetAddr(addr string) {
	s.addr = addr
}

func (s *PeerIdentify) WithContext(c *AppContext) {
	s.AppContext = c
}
