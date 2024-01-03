package tcp

import "time"

type Option func(*option)

type option struct {
	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	HandshakeTimeout time.Duration
}

// specifies the duration for the handshake to complete.Default is 5 second
func WithHandshakeTimeout(duration time.Duration) Option {
	return func(option *option) {
		option.HandshakeTimeout = duration
	}
}

// set write timeout,Default  is no timeout
func WithWriteTimeout(duration time.Duration) Option {
	return func(option *option) {
		option.WriteTimeout = duration
	}
}

// set read timeout,Default  is no timeout
func WithReadTimeout(duration time.Duration) Option {
	return func(option *option) {
		option.ReadTimeout = duration
	}
}
