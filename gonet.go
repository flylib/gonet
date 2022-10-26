package gonet

import (
	"github.com/zjllib/gonet/v3/codec"
	transport2 "github.com/zjllib/gonet/v3/transport"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
)

//一切皆服务
type Service interface {
	// 开启服务
	Start() error
	// 停止服务
	Stop() error
}

var (
	ctx Context //上下文
)

type Hook func(msg *Message)

type Context struct {
	//会话管理
	SessionManager
	//message types
	mMsgTypes map[MessageID]reflect.Type
	//message ids
	mMsgIDs map[reflect.Type]MessageID
	//server types
	sessionType reflect.Type
	//消息编码器
	defaultCodec Codec
	//传输端
	transport transport2.Transport
	//bee worker pool
	workers BeeWorkerPool
	//消息钩子
	mMsgHooks map[MessageID]Hook
}

func (c Context) Start() error {
	return c.transport.Listen()
}

func (c Context) Stop() error {
	return c.transport.Stop()
}

func init() {
	log.SetPrefix("[gonet]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func init() {
	ctx = Context{
		SessionManager: SessionManager{
			pool: sync.Pool{
				New: func() interface{} {
					return reflect.New(ctx.sessionType).Interface()
				},
			},
		},
		mMsgTypes: map[MessageID]reflect.Type{},
		mMsgIDs:   map[reflect.Type]MessageID{},
		mMsgHooks: map[MessageID]Hook{},
	}
}

func NewService(opts ...options) Service {
	option := Option{}
	for _, f := range opts {
		f(&option)
	}
	//传输协议
	ctx.transport = option.transport
	ctx.sessionType = option.transport.SessionType()

	//编码格式
	switch option.contentType {
	case Binary:
		ctx.defaultCodec = codec.BinaryCodec{}
	case Xml:
		ctx.defaultCodec = codec.XmlCodec{}
	case Protobuf:
		ctx.defaultCodec = codec.ProtobufCodec{}
	default:
		ctx.defaultCodec = codec.JsonCodec{}
	}
	cache := option.msgCache
	if cache == nil {
		cache = &MessageList{}
	}
	ctx.workers = createBeeWorkerPool(option.workerPoolSize, cache)
	return ctx
}

//获取会话
func GetSession(id uint64) (Session, bool) {
	value, ok := ctx.sessions.Load(id)
	if ok {
		return value.(Session), ok
	}
	return nil, false
}

//创建会话
func CreateSession() Session {
	obj := ctx.pool.Get()
	ctx.store(atomic.AddUint64(&ctx.incr, 1), obj)
	session := obj.(Session)
	return session
}

//回收会话对象
func RecycleSession(session Session, err error) {
	CacheMsg(&Message{
		Session: session,
		ID:      SessionClose,
		Body:    err,
	})
	//关闭
	session.Close()
	//删除
	ctx.del(session.ID())
	//回收
	ctx.pool.Put(session)
}

//统计会话数量
func SessionCount() int {
	sum := 0
	ctx.sessions.Range(func(key, value interface{}) bool {
		sum++
		return true
	})
	return sum
}

//广播会话
func Broadcast(msg interface{}) {
	ctx.sessions.Range(func(_, item interface{}) bool {
		session, ok := item.(Session)
		if ok {
			session.Send(msg)
		}
		return true
	})
}

//映射消息体
func Route(msgID MessageID, msg interface{}, callback Hook) {
	ctx.Lock()
	defer ctx.Unlock()
	msgType := reflect.TypeOf(msg)
	if _, ok := ctx.mMsgTypes[msgID]; ok {
		panic("error:Duplicate message id")
	}
	if msgType != nil {
		ctx.mMsgIDs[msgType] = msgID
		ctx.mMsgTypes[msgID] = msgType
	}
	if callback != nil {
		ctx.mMsgHooks[msgID] = callback
	}
}

//获取消息ID
func GetMsgID(msg interface{}) (MessageID, bool) {
	msgID, ok := ctx.mMsgIDs[reflect.TypeOf(msg)]
	return msgID, ok
}

//通消息id创建消息体
func CreateMsg(msgID MessageID) interface{} {
	if msg, ok := ctx.mMsgTypes[msgID]; ok {
		return reflect.New(msg).Interface()
	}
	return nil
}

//编码消息
func EncodeMessage(msg interface{}) ([]byte, error) {
	return ctx.defaultCodec.Encode(msg)
}

// 解码消息
func DecodeMessage(msg interface{}, data []byte) error {
	return ctx.defaultCodec.Decode(data, msg)
}

//缓存消息
func CacheMsg(msg *Message) {
	ctx.workers.rcvMsgCh <- msg
}
