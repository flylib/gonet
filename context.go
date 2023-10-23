package gonet

import (
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/goutils/logger/log"
	"reflect"
)

type AppContext struct {
	opt option
	//session manager
	sessionMgr *SessionManager
	//go routine pool
	routines *GoroutinePool
	ILogger
}

func NewContext(options ...Option) *AppContext {
	ctx := &AppContext{
		ILogger: log.NewLogger(),
		opt: option{
			codec:            new(json.Codec),
			netPackageParser: new(DefaultNetPackageParser),
		},
	}

	for _, f := range options {
		f(&ctx.opt)
	}

	if ctx.opt.log != nil {
		ctx.ILogger = ctx.opt.log
	}
	ctx.routines = newGoroutinePool(ctx)
	return ctx
}

// 会话管理
func (c *AppContext) GetSession(id uint64) (ISession, bool) {
	return c.sessionMgr.GetAliveSession(id)
}

func (c *AppContext) InitSessionMgr(sessionType reflect.Type) {
	c.sessionMgr = NewSessionManager(sessionType)
}

func (c *AppContext) CreateSession() ISession {
	idleSession := c.sessionMgr.GetIdleSession()
	idleSession.(ISessionIdentify).ClearIdentify()
	session := idleSession.(ISession)
	c.sessionMgr.AddAliveSession(idleSession)
	c.PushGlobalMessageQueue(newConnectionConnectMessage(session))
	return session
}

func (c *AppContext) RecycleSession(session ISession, err error) {
	c.PushGlobalMessageQueue(newConnectionCloseMessage(session, err))
	session.Close()
	session.(ISessionAbility).ClearAbility()
	c.sessionMgr.RecycleIdleSession(session)
}

func (c *AppContext) SessionCount() int32 {
	return c.sessionMgr.CountAliveSession()
}
func (c *AppContext) Broadcast(msgId uint32, msg any) {
	c.sessionMgr.alive.Range(func(_, item interface{}) bool {
		session, ok := item.(ISession)
		if ok {
			session.Send(msgId, msg)
		}
		return true
	})
}

// message encoding
func (c *AppContext) Marshal(msg any) ([]byte, error) {
	return c.opt.codec.Marshal(msg)
}
func (c *AppContext) Unmarshal(data []byte, v any) error {
	return c.opt.codec.Unmarshal(data, v)
}

// network packet
func (c *AppContext) PackageMessage(s ISession, messageId uint32, v any) ([]byte, error) {
	return c.opt.netPackageParser.Package(s, messageId, v)
}

func (c *AppContext) UnPackageMessage(s ISession, data []byte) (IMessage, int, error) {
	return c.opt.netPackageParser.UnPackage(s, data)
}

// push the message to the routine pool
func (c *AppContext) PushGlobalMessageQueue(msg IMessage) {
	// active defense to avoid too many message
	c.routines.queue <- msg
}
