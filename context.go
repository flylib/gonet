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
	//会话管理
	sessionMgr *SessionManager
	//message types
	mMsgTypes map[MessageID]reflect.Type
	//message ids
	mMsgIDs map[reflect.Type]MessageID
	//server types
	sessionType reflect.Type
	//消息编码器
	defaultCodec codec.Codec
	//传输端
	server IServer
	client IClient
	//bee worker pool
	workers BeeWorkerPool
	//消息钩子
	mMsgHooks map[MessageID]MessageHandler

	name string

	globalLock sync.Mutex

	//包解析器
	IPackageParser
}

func NewContext(opts ...options) *Context {
	option := Option{}
	for _, f := range opts {
		f(&option)
	}
	c := &Context{
		sessionMgr: newSessionManager(option.server.SessionType()),
		mMsgTypes:  map[MessageID]reflect.Type{},
		mMsgIDs:    map[reflect.Type]MessageID{},
		mMsgHooks:  map[MessageID]MessageHandler{},
	}
	//传输协议
	c.server = option.server
	c.server.(IPeerIdentify).WithContext(c)
	c.sessionType = option.server.SessionType()
	if option.serviceName == "" {
		option.serviceName = "gonet"
	}
	c.name = option.serviceName

	//编码格式
	switch option.contentType {
	case codec.Binary:
		c.defaultCodec = binary.BinaryCodec{}
	case codec.Xml:
		c.defaultCodec = xml.XmlCodec{}
	case codec.Protobuf:
		c.defaultCodec = protobuf.ProtobufCodec{}
	default:
		c.defaultCodec = json.JsonCodec{}
	}
	cache := option.msgCache
	if cache == nil {
		cache = &MessageList{}
	}
	c.workers = createBeeWorkerPool(c, option.workerPoolSize, cache)
	c.IPackageParser = &defaultPackageParser{c}
	return c
}

func (c *Context) Name() string {
	return c.name
}

// peer
func (c *Context) Server() IServer {
	return c.server
}
func (c *Context) Client() IClient {
	return c.client
}

// 会话管理
func (c *Context) GetSession(id uint64) (ISession, bool) {
	return c.sessionMgr.getAliveSession(id)
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
		panic("error:Duplicate message id")
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
	return c.defaultCodec.Encode(msg)
}
func (c *Context) DecodeMessage(msg any, data []byte) error {
	return c.defaultCodec.Decode(data, msg)
}

// 缓存消息
func (c *Context) PushGlobalMessageQueue(session ISession, msg IMessage) {
	c.workers.rcvMsgCh <- event{session: session, message: msg}
}
