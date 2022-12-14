package gonet

import (
	"github.com/zjllib/gonet/v3/codec"
	"github.com/zjllib/gonet/v3/codec/binary"
	"github.com/zjllib/gonet/v3/codec/json"
	"github.com/zjllib/gonet/v3/codec/protobuf"
	"github.com/zjllib/gonet/v3/codec/xml"
	"github.com/zjllib/gonet/v3/transport"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	goNetContext goNet //上下文
)

func init() {
	log.SetPrefix("[gonet]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func init() {
	goNetContext = goNet{
		SessionManager: SessionManager{
			pool: sync.Pool{
				New: func() interface{} {
					return reflect.New(goNetContext.sessionType).Interface()
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
	defaultCodec codec.Codec
	//传输端
	server transport.IServer
	//bee worker pool
	workers BeeWorkerPool
	//消息钩子
	mMsgHooks map[MessageID]Hook

	name string
}

func (c *goNet) Name() string {
	return c.name
}

func (c *goNet) Start() error {
	return c.server.Listen()
}

func (c *goNet) Stop() error {
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
	session.(interface{ SetID(id uint64) }).SetID(id)
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
	goNetContext.server = option.server
	goNetContext.sessionType = option.server.SessionType()
	if option.serviceName == "" {
		option.serviceName = "goNet"
	}
	goNetContext.name = option.serviceName

	//编码格式
	switch option.contentType {
	case codec.Binary:
		goNetContext.defaultCodec = binary.BinaryCodec{}
	case codec.Xml:
		goNetContext.defaultCodec = xml.XmlCodec{}
	case codec.Protobuf:
		goNetContext.defaultCodec = protobuf.ProtobufCodec{}
	default:
		goNetContext.defaultCodec = json.JsonCodec{}
	}
	cache := option.msgCache
	if cache == nil {
		cache = &MessageList{}
	}
	goNetContext.workers = createBeeWorkerPool(option.workerPoolSize, cache)
	return &goNetContext
}

//获取会话
func GetSession(id uint64) (transport.ISession, bool) {
	value, ok := goNetContext.sessions.Load(id)
	if ok {
		return value.(transport.ISession), ok
	}
	return nil, false
}

//创建会话
func CreateSession() transport.ISession {
	obj := goNetContext.pool.Get()
	goNetContext.store(atomic.AddUint64(&goNetContext.incr, 1), obj)
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
	goNetContext.del(session.ID())
	//回收
	goNetContext.pool.Put(session)
}

//统计会话数量
func SessionCount() int {
	sum := 0
	goNetContext.sessions.Range(func(key, value interface{}) bool {
		sum++
		return true
	})
	return sum
}

//广播会话
func Broadcast(msg interface{}) {
	goNetContext.sessions.Range(func(_, item interface{}) bool {
		session, ok := item.(transport.ISession)
		if ok {
			session.Send(msg)
		}
		return true
	})
}

//映射消息体
func Route(msgID MessageID, msg interface{}, callback Hook) {
	goNetContext.Lock()
	defer goNetContext.Unlock()
	msgType := reflect.TypeOf(msg)
	if _, ok := goNetContext.mMsgTypes[msgID]; ok {
		panic("error:Duplicate message id")
	}
	if msgType != nil {
		goNetContext.mMsgIDs[msgType] = msgID
		goNetContext.mMsgTypes[msgID] = msgType
	}
	if callback != nil {
		goNetContext.mMsgHooks[msgID] = callback
	}
}

//获取消息ID
func GetMsgID(msg interface{}) (MessageID, bool) {
	msgID, ok := goNetContext.mMsgIDs[reflect.TypeOf(msg)]
	return msgID, ok
}

//通消息id创建消息体
func CreateMsg(msgID MessageID) interface{} {
	if msg, ok := goNetContext.mMsgTypes[msgID]; ok {
		return reflect.New(msg).Interface()
	}
	return nil
}

//编码消息
func EncodeMessage(msg interface{}) ([]byte, error) {
	return goNetContext.defaultCodec.Encode(msg)
}

// 解码消息
func DecodeMessage(msg interface{}, data []byte) error {
	return goNetContext.defaultCodec.Decode(data, msg)
}

//缓存消息
func CacheMessage(session transport.ISession, msg *Message) {
	msg.Head.setSession(session)
	goNetContext.workers.rcvMsgCh <- msg
}
