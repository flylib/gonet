package gonet

import (
	"github.com/zjllib/gonet/v3/codec/json"
	"reflect"
	"sync"
)

type Context struct {
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
	//bee worker pool
	bees              BeeWorkerPool
	maxWorkerPoolSize int
	//cache for messages
	msgCache IMessageCache

	globalLock sync.Mutex

	//包解析器
	INetPackageParser
	//0意味着无限制
	maxSessionCount int
	//contentType support json/xml/binary/protobuf
	contentType string
}

func NewContext(handlers ...Option) *Context {
	ctx := &Context{
		mMsgTypes: make(map[MessageID]reflect.Type),
		mMsgIDs:   make(map[reflect.Type]MessageID),
		mMsgHooks: make(map[MessageID]MessageHandler),
	}
	for _, handler := range handlers {
		err := handler(ctx)
		if err != nil {
			panic(err)
		}
	}
	//编码格式
	if ctx.codec == nil {
		ctx.codec = new(json.Codec)
	}
	if ctx.msgCache == nil {
		ctx.msgCache = new(DefaultMessageCacheList)
	}
	ctx.bees = createBeeWorkerPool(ctx, ctx.maxWorkerPoolSize, ctx.msgCache)
	ctx.INetPackageParser = &DefaultNetPackageParser{ctx}
	return ctx
}

// 会话管理
func (c *Context) GetSession(id uint64) (ISession, bool) {
	return c.sessionMgr.getAliveSession(id)
}

func (c *Context) InitSessionMgr(sessionType reflect.Type) {
	c.sessionMgr = newSessionManager(sessionType)
}

func (c *Context) CreateSession() ISession {
	idleSession := c.sessionMgr.getIdleSession()
	idleSession.(ISessionIdentify).ClearIdentify()
	session := idleSession.(ISession)
	c.sessionMgr.addAliveSession(idleSession)
	session.(ISessionAbility).InitSendChanel()
	c.PushGlobalMessageQueue(newSessionConnectMessage(session))
	return session
}
func (c *Context) RecycleSession(session ISession, err error) {
	c.PushGlobalMessageQueue(newSessionCloseMessage(session, err))
	session.Close()
	session.(ISessionAbility).StopAbility()
	c.sessionMgr.recycleIdleSession(session)
}
func (c *Context) SessionCount() int {
	return int(c.sessionMgr.CountAliveSession())
}
func (c *Context) Broadcast(msg interface{}) {
	c.sessionMgr.alive.Range(func(_, item interface{}) bool {
		session, ok := item.(ISession)
		if ok {
			session.Send(msg)
		}
		return true
	})
}

// 消息管理
func (c *Context) Route(msgID MessageID, msg any, callback MessageHandler) {
	c.globalLock.Lock()
	defer c.globalLock.Unlock()

	msgType := reflect.TypeOf(msg)
	if _, ok := c.mMsgTypes[msgID]; ok {
		panic("Duplicate message id")
	}
	if msgType != nil {
		c.mMsgIDs[msgType] = msgID
		c.mMsgTypes[msgID] = msgType
	}
	if callback != nil {
		c.mMsgHooks[msgID] = callback
	}
}

func (c *Context) GetMsgID(msg interface{}) (MessageID, bool) {
	msgID, ok := c.mMsgIDs[reflect.TypeOf(msg)]
	return msgID, ok
}
func (c *Context) CreateMsg(msgID MessageID) interface{} {
	if msg, ok := c.mMsgTypes[msgID]; ok {
		return reflect.New(msg).Interface()
	}
	return nil
}

// 消息编码
func (c *Context) EncodeMessage(msg any) ([]byte, error) {
	return c.codec.Encode(msg)
}
func (c *Context) DecodeMessage(msg any, data []byte) error {
	return c.codec.Decode(data, msg)
}

// 缓存消息
func (c *Context) PushGlobalMessageQueue(msg IMessage) {
	//主动防御，避免消息过多
	c.bees.rcvMsgCh <- msg
}
