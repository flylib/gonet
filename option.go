package goNet

import (
	"time"
)

var Opts = &Options{}

// Option represents the optional function.
type Option func(opts *Options)

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
}

// WithOptions accepts the whole options config.
func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

//set addr
func WithAddr(addr string) Option {
	return func(opts *Options) {
		opts.Addr = addr
	}
}

//set peer type
func WithPeerType(peerType PeerType) Option {
	return func(opts *Options) {
		opts.PeerType = peerType
	}
}

//set read  duration
func WithReadDeadline(dur time.Duration) Option {
	return func(opts *Options) {
		opts.ReadDeadline = dur
	}
}

//set write  duration
func WithWriteDeadline(dur time.Duration) Option {
	return func(opts *Options) {
		opts.WriteDeadline = dur
	}
}

//bind addr
func WithRoutinePoolSize(size int) Option {
	return func(opts *Options) {
		opts.PoolSize = size
	}
}

// WithPanicHandler sets up panic handler.
func WithPanicHandler(panicHandler func(interface{})) Option {
	return func(opts *Options) {
		opts.PanicHandler = panicHandler
	}
}

//通讯协议
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
