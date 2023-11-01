package gonet

import (
	"github.com/flylib/interface/codec"
	ilog "github.com/flylib/interface/log"
	"reflect"
)

type Context struct {
	//session manager
	sessionMgr *SessionManager
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
}

func NewContext(options ...Option) *Context {
	ctx := &Context{
		INetPackager: &DefaultNetPackager{},
	}

	for _, f := range options {
		f(ctx)
	}

	ctx.routines = newGoroutinePool(ctx)
	return ctx
}

// 会话管理
func (c *Context) GetSession(id uint64) (ISession, bool) {
	return c.sessionMgr.GetAliveSession(id)
}

func (c *Context) InitSessionMgr(sessionType reflect.Type) {
	c.sessionMgr = NewSessionManager(sessionType)
}

func (c *Context) CreateSession() ISession {
	idleSession := c.sessionMgr.GetIdleSession()
	idleSession.(ISessionIdentify).ClearIdentify()
	session := idleSession.(ISession)
	c.sessionMgr.AddAliveSession(idleSession)
	c.PushGlobalMessageQueue(newConnectionConnectMessage(session))
	return session
}

func (c *Context) RecycleSession(session ISession, err error) {
	c.PushGlobalMessageQueue(newConnectionCloseMessage(session, err))
	session.Close()
	session.(ISessionAbility).ClearAbility()
	c.sessionMgr.RecycleIdleSession(session)
}

func (c *Context) SessionCount() int32 {
	return c.sessionMgr.CountAliveSession()
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
