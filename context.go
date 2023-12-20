package gonet

import (
	"github.com/flylib/interface/codec"
	ilog "github.com/flylib/interface/log"
	"reflect"
)

type TransportProtocol string

const (
	TCP  TransportProtocol = "tcp"
	KCP  TransportProtocol = "kcp"
	UDP  TransportProtocol = "udp"
	WS   TransportProtocol = "websocket"
	HTTP TransportProtocol = "http"
	QUIC TransportProtocol = "quic"
	RPC  TransportProtocol = "rpc"
)

type invalidData struct {
}

var zeroData = invalidData{}

type Context struct {
	//session manager
	sessionMgr *sessionManager
	//go routine pool
	routines *GoroutinePool

	//Message callback processing
	messageHandler  MessageHandler
	maxSessionCount int
	//routine pool config
	poolCfg poolConfig
	//message codec
	codec.ICodec
	ilog.ILogger
	//net package parser
	INetPackager

	sessionType reflect.Type
}

func NewContext(options ...Option) *Context {
	ctx := &Context{
		INetPackager: &DefaultNetPackager{},
	}

	for _, f := range options {
		f(ctx)
	}

	if ctx.ICodec == nil {
		panic("nil ICodec")
	}

	if ctx.ILogger == nil {
		panic("nil ILogger")
	}

	if ctx.sessionType == nil {
		panic("nil sessionType")
	}

	ctx.routines = newGoroutinePool(ctx)
	ctx.sessionMgr = newSessionManager(ctx.sessionType)
	return ctx
}

// 会话管理
func (c *Context) GetSession(id uint64) (ISession, bool) {
	return c.sessionMgr.getAliveSession(id)
}
func (c *Context) CreateSession() ISession {
	idleSession := c.sessionMgr.getIdleSession()
	idleSession.(interface{ Clear() }).Clear()
	session := idleSession.(ISession)
	c.sessionMgr.addAliveSession(idleSession)
	c.PushGlobalMessageQueue(newConnectionConnectMessage(session))
	return session
}
func (c *Context) RecycleSession(session ISession, err error) {
	c.PushGlobalMessageQueue(newConnectionCloseMessage(session, err))
	session.Close()
	session.(interface{ Clear() }).Clear()
	c.sessionMgr.recycleIdleSession(session)
}

func (c *Context) SessionCount() int32 {
	return c.sessionMgr.countAliveSession()
}

func (c *Context) Broadcast(msgId uint32, msg any) {
	c.sessionMgr.alive.Range(func(_, item interface{}) bool {
		session, ok := item.(ISession)
		if ok {
			session.Send(msgId, msg)
		}
		return true
	})
}

// push the message to the routine pool
func (c *Context) PushGlobalMessageQueue(msg IMessage) {
	// active defense to avoid too many message
	c.routines.queue <- msg
}
