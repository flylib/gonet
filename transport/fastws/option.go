package fastws

import (
	"time"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

type Option func(*option)

type option struct {
	websocket.FastHTTPUpgrader
	HandshakeTimeout time.Duration
}

func WithHandshakeTimeout(t time.Duration) Option {
	return func(o *option) {
		o.HandshakeTimeout = t
	}
}

func WithReadBufferSize(size int) Option {
	return func(o *option) {
		o.ReadBufferSize = size
	}
}

func WithWriteBufferSize(size int) Option {
	return func(o *option) {
		o.WriteBufferSize = size
	}
}

func WithEnableCompression(enable bool) Option {
	return func(o *option) {
		o.EnableCompression = enable
	}
}

func WithCheckOrigin(fn func(ctx *fasthttp.RequestCtx) bool) Option {
	return func(o *option) {
		o.CheckOrigin = fn
	}
}
