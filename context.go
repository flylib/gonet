package gonet

import (
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/goutils/logger/log"
	"reflect"
	"sync"
)

type AppContext struct {
	//session manager
	sessionMgr *SessionManager
	//message types
	mMsgTypes map[MessageID]reflect.Type
	//message ids
	mMsgIDs map[reflect.Type]MessageID
	//msg route handler
	mMsgHooks map[MessageID]MessageHandler

	//message codec
	codec ICodec

	globalLock sync.Mutex

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
		mMsgTypes: map[MessageID]reflect.Type{
			MessageID_Invalid: nil,
		},
		mMsgIDs:          make(map[reflect.Type]MessageID),
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
func (c *AppContext) Route(msgID MessageID, msg any, callback MessageHandler) {
	c.globalLock.Lock()
	defer c.globalLock.Unlock()

	msgType := reflect.TypeOf(msg)
	if _, ok := c.mMsgTypes[msgID]; ok {
		panic("Duplicate message")
	}
	if msgType != nil {
		c.mMsgIDs[msgType] = msgID
		c.mMsgTypes[msgID] = msgType
	}
	if callback != nil {
		c.mMsgHooks[msgID] = callback
	}
}
func (c *AppContext) GetMsgID(msg interface{}) (MessageID, bool) {
	msgID, ok := c.mMsgIDs[reflect.TypeOf(msg)]
	return msgID, ok
}
func (c *AppContext) CreateMsg(msgID MessageID) interface{} {
	if msg, ok := c.mMsgTypes[msgID]; ok {
		return reflect.New(msg).Interface()
	}
	return nil
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
func (c *AppContext) PackageMessage(msg any) ([]byte, error) {
	return c.netPackageParser.Package(c, msg)
}
func (c *AppContext) UnPackageMessage(s ISession, data []byte) (IMessage, int, error) {
	return c.netPackageParser.UnPackage(c, s, data)
}

// 缓存消息
func (c *AppContext) PushGlobalMessageQueue(msg IMessage) {
	//todo 主动防御，避免消息过多
	c.workers.queue <- msg
}
