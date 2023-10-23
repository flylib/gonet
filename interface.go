package gonet

// Transport layer related
type (
	//server
	IServer interface {
		Addr() string
		Listen(addr string) error
		Stop() error
	}
	//client
	IClient interface {
		Dial(addr string) (ISession, error)
	}
)
