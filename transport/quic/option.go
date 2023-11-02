package quic

import "time"

type Option func(*option)

type option struct {
	//specifies the duration for the handshake to complete.Default is 5 second
	HandshakeTimeout time.Duration
}

func WithHandshakeTimeout(t time.Duration) Option {
	return func(option *option) {
		option.HandshakeTimeout = t
	}
}
