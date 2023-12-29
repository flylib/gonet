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

func SetContext(options ...Option) *Context {
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
	return ctx
}

func DefaultContext() *Context {
	return defaultCtx
}

// session manager
func (c *Context) GetSessionManager() *SessionManager {
	return c.sessionManager
}

// event handler
func (c *Context) GetEventHandler() IEventHandler {
	return c.eventHandler
}

// net packager
func (c *Context) GetNetPackager() INetPackager {
	return c.netPackager
}

// net packager
func (c *Context) GetCodec() codec.ICodec {
	return c.codec
}

// async runtime
func (c *Context) GetAsyncRuntime() *AsyncRuntime {
	return c.asyncRuntime
}
