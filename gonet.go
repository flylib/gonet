package gonet

import (
	"github.com/zjllib/gonet/v3/codec"
	"github.com/zjllib/gonet/v3/transport"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	ctx goNet //上下文
)

func init() {
	log.SetPrefix("[gonet]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func init() {
	ctx = goNet{
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

type Hook func(msg *Message)

type goNet struct {
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
	server transport.IServer
	//bee worker pool
	workers BeeWorkerPool
	//消息钩子
	mMsgHooks map[MessageID]Hook

	name string
}

func (c goNet) Name() string {
	return c.name
}

func (c goNet) Start() error {
	return c.server.Listen()
}

func (c goNet) Stop() error {
	return c.server.Stop()
}

//会话管理
type SessionManager struct {
	sync.RWMutex
	incr     uint64    //流水号
	sessions sync.Map  //所有链接
	pool     sync.Pool //临时对象池
}

func (s *SessionManager) store(id uint64, session interface{}) {
	session.(interface{ setID(id uint64) }).setID(id)
	s.sessions.Store(id, session)
}

func (s *SessionManager) del(id uint64) {
	s.sessions.Delete(id)
}

func NewService(opts ...options) IService {
	option := Option{}
	for _, f := range opts {
		f(&option)
	}
	//传输协议
	ctx.server = option.server
	ctx.sessionType = option.server.SessionType()
	if option.serviceName == "" {
		option.serviceName = "goNet"
	}
	ctx.name = option.serviceName

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
func GetSession(id uint64) (transport.ISession, bool) {
	value, ok := ctx.sessions.Load(id)
	if ok {
		return value.(transport.ISession), ok
	}
	return nil, false
}

//创建会话
func CreateSession() transport.ISession {
	obj := ctx.pool.Get()
	ctx.store(atomic.AddUint64(&ctx.incr, 1), obj)
	session := obj.(transport.ISession)
	return session
}

//回收会话对象
func RecycleSession(session transport.ISession, err error) {
	CacheMessage(session, &Message{
		ID:   SessionClose,
		Body: err,
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
		session, ok := item.(transport.ISession)
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
func CacheMessage(session transport.ISession, msg *Message) {
	msg.Head.setSession(session)
	ctx.workers.rcvMsgCh <- msg
}
