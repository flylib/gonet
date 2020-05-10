package goNet

import (
	"time"
)

type ProtoCol string

const (
	TCP  ProtoCol = "tcp"
	KCP  ProtoCol = "kcp"
	UDP  ProtoCol = "udp"
	WS   ProtoCol = "websocket"
	HTTP ProtoCol = "http"
	QUIC ProtoCol = "quic"
	RPC  ProtoCol = "rpc"
)

//options
type Options struct {
	//listen or dial addr
	Addr string
	//peer type
	PeerType PeerType
	//SetWriteDeadline sets the write deadline or read deadline on the underlying connection.
	ReadDeadline, WriteDeadline time.Duration
	//set the routine pool size
	//0.mean use default set
	PoolSize int
	// PanicHandler is used to handle panics from each worker goroutine.
	PanicHandler func(interface{})
	//Maximum number of connections allowed
	//0.mean no limit
	AllowMaxConn int
}
