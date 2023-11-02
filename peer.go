package gonet

// Transport layer related
type (
	//server
	IServer interface {
		Listen(addr string) error
		Close() error
		Addr() string
	}
	//client
	IClient interface {
		Dial(addr string) (ISession, error)
		Close() error
	}
)

type IPeerIdentify interface {
	Addr() string
	SetAddr(addr string)
	WithContext(c *Context)
}
type PeerIdentify struct {
	*Context
	uuid string
	addr string
}

func (s *PeerIdentify) Addr() string {
	return s.addr
}

func (s *PeerIdentify) SetAddr(addr string) {
	s.addr = addr
}

func (s *PeerIdentify) WithContext(c *Context) {
	s.Context = c
}
