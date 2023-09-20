package gonet

import (
	"github.com/zjllib/gonet/v3/codec"
	"github.com/zjllib/gonet/v3/codec/binary"
	"github.com/zjllib/gonet/v3/codec/json"
	"github.com/zjllib/gonet/v3/codec/protobuf"
	"github.com/zjllib/gonet/v3/codec/xml"
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
	codec codec.ICodec
	//bee worker pool
	workers           BeeWorkerPool
	maxWorkerPoolSize int
	//cache for messages
	msgCache IEventCache

	globalLock sync.Mutex

	//包解析器
	IPackageParser
	//0意味着无限制
	maxSessionCount int
	//contentType support json/xml/binary/protobuf
	contentType string
}

func NewContext(handlers ...HandlerFunc) *Context {
	ctx := &Context{
		mMsgTypes: map[MessageID]reflect.Type{},
		mMsgIDs:   map[reflect.Type]MessageID{},
		mMsgHooks: map[MessageID]MessageHandler{},
	}
	for _, handler := range handlers {
		err := handler(ctx)
		if err != nil {
			panic(err)
		}
	}
	//编码格式
	switch ctx.contentType {
	case codec.Binary:
		ctx.codec = binary.BinaryCodec{}
	case codec.Xml:
		ctx.codec = xml.XmlCodec{}
	case codec.Protobuf:
		ctx.codec = protobuf.ProtobufCodec{}
	default:
		ctx.codec = json.JsonCodec{}
	}
	cache := ctx.msgCache
	if cache == nil {
		cache = &MessageList{}
	}
	ctx.workers = createBeeWorkerPool(ctx, ctx.maxWorkerPoolSize, cache)
	ctx.IPackageParser = &defaultPackageParser{ctx}
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
	c.PushGlobalMessageQueue(session, NewSessionMessage)
	return session
}
func (c *Context) RecycleSession(session ISession, err error) {
	c.PushGlobalMessageQueue(session, &Message{
		id:   SessionClose,
		body: err,
	})
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
func (c *Context) PushGlobalMessageQueue(session ISession, msg IMessage) {
	c.workers.rcvMsgCh <- event{session: session, message: msg}
}
