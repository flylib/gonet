package gnet

import (
	"github.com/panjf2000/gnet/v2"
	"time"
)

// LoadBalancing represents the type of load-balancing algorithm.
type LoadBalancing int

const (
	// RoundRobin assigns the next accepted connection to the event-loop by polling event-loop list.
	RoundRobin LoadBalancing = iota

	// LeastConnections assigns the next accepted connection to the event-loop that is
	// serving the least number of active connections at the current time.
	LeastConnections

	// SourceAddrHash assigns the next accepted connection to the event-loop by hashing the remote address.
	SourceAddrHash
)

// TCPSocketOpt is the type of TCP socket options.
type TCPSocketOpt int

// Available TCP socket options.
const (
	TCPNoDelay TCPSocketOpt = iota
	TCPDelay
)

type Option func(*option)

type option struct {
	gnet.Options
}

// WithMulticore sets up multi-cores in gnet engine.
func WithMulticore(multicore bool) Option {
	return func(opts *option) {
		opts.Multicore = multicore
	}
}

// WithLockOSThread sets up LockOSThread mode for I/O event-loops.
func WithLockOSThread(lockOSThread bool) Option {
	return func(opts *option) {
		opts.LockOSThread = lockOSThread
	}
}

// WithReadBufferCap sets up ReadBufferCap for reading bytes.
func WithReadBufferCap(readBufferCap int) Option {
	return func(opts *option) {
		opts.ReadBufferCap = readBufferCap
	}
}

// WithWriteBufferCap sets up WriteBufferCap for pending bytes.
func WithWriteBufferCap(writeBufferCap int) Option {
	return func(opts *option) {
		opts.WriteBufferCap = writeBufferCap
	}
}

// WithLoadBalancing sets up the load-balancing algorithm in gnet engine.
func WithLoadBalancing(lb LoadBalancing) Option {
	return func(opts *option) {
		opts.LB = gnet.LoadBalancing(lb)
	}
}

// WithNumEventLoop sets up NumEventLoop in gnet engine.
func WithNumEventLoop(numEventLoop int) Option {
	return func(opts *option) {
		opts.NumEventLoop = numEventLoop
	}
}

// WithReusePort sets up SO_REUSEPORT socket option.
func WithReusePort(reusePort bool) Option {
	return func(opts *option) {
		opts.ReusePort = reusePort
	}
}

// WithReuseAddr sets up SO_REUSEADDR socket option.
func WithReuseAddr(reuseAddr bool) Option {
	return func(opts *option) {
		opts.ReuseAddr = reuseAddr
	}
}

// WithTCPKeepAlive sets up the SO_KEEPALIVE socket option with duration.
func WithTCPKeepAlive(tcpKeepAlive time.Duration) Option {
	return func(opts *option) {
		opts.TCPKeepAlive = tcpKeepAlive
	}
}

// WithTCPNoDelay enable/disable the TCP_NODELAY socket option.
func WithTCPNoDelay(tcpNoDelay TCPSocketOpt) Option {
	return func(opts *option) {
		opts.TCPNoDelay = gnet.TCPSocketOpt(tcpNoDelay)
	}
}

// WithSocketRecvBuffer sets the maximum socket receive buffer in bytes.
func WithSocketRecvBuffer(recvBuf int) Option {
	return func(opts *option) {
		opts.SocketRecvBuffer = recvBuf
	}
}

// WithSocketSendBuffer sets the maximum socket send buffer in bytes.
func WithSocketSendBuffer(sendBuf int) Option {
	return func(opts *option) {
		opts.SocketSendBuffer = sendBuf
	}
}

// WithTicker indicates that a ticker is set.
func WithTicker(ticker bool) Option {
	return func(opts *option) {
		opts.Ticker = ticker
	}
}

// WithMulticastInterfaceIndex sets the interface name where UDP multicast sockets will be bound to.
func WithMulticastInterfaceIndex(idx int) Option {
	return func(opts *option) {
		opts.MulticastInterfaceIndex = idx
	}
}
