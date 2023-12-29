package gonet

// Transport layer related
type (
	//server interface
	IServer interface {
		Listen(addr string) error
		Close() error
		Addr() string
	}
	//client interface
	IClient interface {
		Dial(addr string) (ISession, error)
		Close() error
	}
)

type PeerCommon struct {
	addr string
}

func (s *PeerCommon) Addr() string {
	return s.addr
}

func (s *PeerCommon) SetAddr(addr string) {
	s.addr = addr
}
