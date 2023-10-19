package gonet

import "net"

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
	//session
	ISession interface {
		//ID
		ID() uint64
		//close the connection
		Close() error
		//send the message to the other side
		Send(msgID uint32, msg any) error
		//remote addr
		RemoteAddr() net.Addr
		//convenient session storage data
		Store(value any)
		//load the data
		Load() (value any)
	}
)
