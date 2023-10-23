package gonet

import (
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/goutils/logger/log"
	"reflect"
)

type AppContext struct {
	callback MessageHandler
	//session manager
	sessionMgr *SessionManager

	//message codec
	codec ICodec

	//net package parser
	netPackageParser INetPackageParser
	//0意味着无限制
	maxSessionCount int

	workers       *GoroutinePool
	workerOptions []goroutinePoolOption
	ILogger
}

func NewContext(options ...Option) *AppContext {
	ctx := &AppContext{
		codec:            new(json.Codec),
		ILogger:          log.NewLogger(),
		netPackageParser: new(DefaultNetPackageParser),
	}
	for _, f := range options {
		err := f(ctx)
		if err != nil {
			panic(err)
		}
	}
	ctx.workers = newGoroutinePool(ctx, ctx.workerOptions...)
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
	c.PushGlobalMessageQueue(newConnectionConnectMessage(idleSession))
	return session
}

func (c *AppContext) RecycleSession(session ISession, err error) {
	c.PushGlobalMessageQueue(newConnectionCloseMessage(session, err))
	session.Close()
	session.(ISessionAbility).ClearAbility()
	c.sessionMgr.RecycleIdleSession(session)
}

func (c *AppContext) SessionCount() uint32 {
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
	return c.codec.Marshal(msg)
}
func (c *AppContext) Unmarshal(data []byte, v any) error {
	return c.codec.Unmarshal(data, v)
}

// network packet
func (c *AppContext) PackageMessage(messageId uint32, v any) ([]byte, error) {
	return c.netPackageParser.Package(messageId, v)
}

func (c *AppContext) UnPackageMessage(s ISession, data []byte) (IMessage, int, error) {
	return c.netPackageParser.UnPackage(s, data)
}

// push the message to the routine pool
func (c *AppContext) PushGlobalMessageQueue(msg IMessage) {
	// active defense to avoid too many message
	c.workers.queue <- msg
}
