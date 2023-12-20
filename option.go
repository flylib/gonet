package gonet

import (
	"github.com/flylib/interface/codec"
	ilog "github.com/flylib/interface/log"
	"reflect"
)

type Option func(o *Context)

// Default:0 means is no limit
func WithMessageHandler(handler MessageHandler) Option {
	return func(o *Context) {
		o.messageHandler = handler
	}
}

// Default:0 means is no limit
func WithMaxSessions(max int) Option {
	return func(o *Context) {
		o.maxSessionCount = max
	}
}

// Default is runtime.NumCPU(), means no goroutines will be dynamically scaled
func WithPoolMaxRoutines(num int32) Option {
	return func(o *Context) {

		o.poolCfg.maxNum = num
	}
}

// allow max idle routines
func WithPoolMaxIdleRoutines(num int32) Option {
	return func(o *Context) {
		o.poolCfg.maxIdleNum = num
	}
}

// Default 512,global queue size
func WithGQSize(size int32) Option {
	return func(o *Context) {
		o.poolCfg.queueSize = size
	}
}

// network package paser
func WithNetPackager(packager INetPackager) Option {
	return func(o *Context) {
		o.INetPackager = packager
	}
}

// Default json codec, message codec
func MustWithCodec(codec codec.ICodec) Option {
	return func(o *Context) {
		o.ICodec = codec
	}
}

// set logger
func MustWithLogger(l ilog.ILogger) Option {
	return func(o *Context) {
		o.ILogger = l
	}
}

// set SessionType
func MustWithSessionType(t reflect.Type) Option {
	return func(o *Context) {
		o.sessionType = t
	}
}
