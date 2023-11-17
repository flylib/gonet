package ws

import (
	"github.com/gorilla/websocket"
	"time"
)

type Option func(*option)

type option struct {
	//specifies the duration for the handshake to complete.Default is 5 second
	websocket.Upgrader //websocket升级器
}

func WithHandshakeTimeout(t time.Duration) Option {
	return func(option *option) {
		option.HandshakeTimeout = t
	}
}

func WithReadBufferSize(size int) Option {
	return func(option *option) {
		option.ReadBufferSize = size
	}
}

func WithWriteBufferSize(size int) Option {
	return func(option *option) {
		option.WriteBufferSize = size
	}
}

func WithEnableCompression(enable bool) Option {
	return func(option *option) {
		option.EnableCompression = enable
	}
}
