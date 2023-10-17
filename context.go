package gonet

import (
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/goutils/logger/log"
	"reflect"
)

type AppContext struct {
	//session manager
	sessionMgr *SessionManager

	//msg route handler
	mMsgHooks map[MessageID]MessageHandler

	//message codec
	codec ICodec

	//net package parser
	netPackageParser INetPackageParser
	//0意味着无限制
	maxSessionCount int
	//contentType support json/xml/binary/protobuf
	contentType string

	workers       *GoroutinePool
	workerOptions []goroutinePoolOption
	ILogger
}

func NewContext(options ...Option) *AppContext {
	ctx := &AppContext{
		mMsgHooks:        make(map[MessageID]MessageHandler),
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
	c.sessionMgr = newSessionManager(sessionType)
}

func (c *AppContext) CreateSession() ISession {
	idleSession := c.sessionMgr.GetIdleSession()
	idleSession.(ISessionIdentify).ClearIdentify()
	session := idleSession.(ISession)
	c.sessionMgr.AddAliveSession(idleSession)
	session.(ISessionAbility).InitSendChanel()
	c.PushGlobalMessageQueue(newSessionConnectMessage(session))
	return session
}

func (c *AppContext) RecycleSession(session ISession, err error) {
	c.PushGlobalMessageQueue(newSessionCloseMessage(session, err))
	session.Close()
	session.(ISessionAbility).StopAbility()
	c.sessionMgr.RecycleIdleSession(session)
}
func (c *AppContext) SessionCount() int {
	return int(c.sessionMgr.CountAliveSession())
}
func (c *AppContext) Broadcast(msg interface{}) {
	c.sessionMgr.alive.Range(func(_, item interface{}) bool {
		session, ok := item.(ISession)
		if ok {
			session.Send(msg)
		}
		return true
	})
}

// 消息管理
func (c *AppContext) Route(msgID MessageID, callback MessageHandler) {
	if _, ok := c.mMsgHooks[msgID]; ok {
		panic("Duplicate message")
	}
	if callback != nil {
		c.mMsgHooks[msgID] = callback
	}
}

func (c *AppContext) GetMessageHandler(msgID MessageID) (MessageHandler, bool) {
	f, ok := c.mMsgHooks[msgID]
	return f, ok
}

// 消息编码
func (c *AppContext) EncodeMessage(msg any) ([]byte, error) {
	return c.codec.Marshal(msg)
}
func (c *AppContext) DecodeMessage(msg any, data []byte) error {
	return c.codec.Unmarshal(data, msg)
}
func (c *AppContext) PackageMessage(messageId MessageID, v any) ([]byte, error) {
	return c.netPackageParser.Package(messageId, v)
}

func (c *AppContext) UnPackageMessage(data []byte) (IMessage, int, error) {
	return c.netPackageParser.UnPackage(data)
}

// 缓存消息
func (c *AppContext) PushGlobalMessageQueue(msg IMessage) {
	//todo 主动防御，避免消息过多
	c.workers.queue <- msg
}
