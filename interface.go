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
		Send(msg any) error
		//remote addr
		RemoteAddr() net.Addr
		//convenient session storage data
		Store(value any)
		//load the data
		Load() (value any, ok bool)
	}
	ISessionIdentify interface {
		ID() uint64
		ClearIdentify()
		SetID(id uint64)
		UpdateID(id uint64)
		WithContext(c *AppContext)
		IsClosed() bool
		SetClosedStatus()
	}
	ISessionAbility interface {
		Store(val any)
		Load() (val any, ok bool)
		InitSendChanel()
		PushSendChannel(buf []byte)
		SendLoop(handler func([]byte))
		StopAbility()
	}
	IPeerIdentify interface {
		Addr() string
		SetAddr(addr string)
		WithContext(c *AppContext)
	}
)

// message related
type (
	IMessage interface {
		ID() MessageID
		Payload() []byte
	}
	IMessageCache interface {
		Size() int
		Push(IMessage)
		Pop() IMessage
	}
)
