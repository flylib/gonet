package gonet

import (
	"github.com/flylib/interface/codec"
	ilog "github.com/flylib/interface/log"
	"reflect"
)

var (
	defaultCtx *Context
)

type Context struct {
	//session manager
	sessionManager *SessionManager
	//asyncRuntime
	asyncRuntime *AsyncRuntime

	//Message callback processing
	eventHandler    IEventHandler
	maxSessionCount int
	//routine pool config
	poolCfg RuntimeConfig
	//message codec
	codec codec.ICodec
	ilog.ILogger
	//net package parser
	netPackager INetPackager

	sessionType reflect.Type
}

func SetupContext(options ...Option) *Context {
	ctx := &Context{
		netPackager: &DefaultNetPackager{},
	}

	for _, f := range options {
		f(ctx)
	}

	if ctx.codec == nil {
		panic("nil ICodec")
	}

	if ctx.ILogger == nil {
		panic("nil ILogger")
	}

	if ctx.sessionType == nil {
		panic("nil sessionType")
	}

	ctx.asyncRuntime = newAsyncRuntime(ctx)
	ctx.sessionManager = NewSessionManager(ctx.sessionType)
	defaultCtx = ctx
	return ctx
}

func DefaultContext() *Context {
	return defaultCtx
}

// session manager
func GetSessionManager() *SessionManager {
	return defaultCtx.sessionManager
}

// event handler
func GetEventHandler() IEventHandler {
	return defaultCtx.eventHandler
}

// net packager
func GetNetPackager() INetPackager {
	return defaultCtx.netPackager
}

// net packager
func GetCodec() codec.ICodec {
	return defaultCtx.codec
}

// async runtime
func GetAsyncRuntime() *AsyncRuntime {
	return defaultCtx.asyncRuntime
}
