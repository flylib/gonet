package gonet

import (
	"github.com/zjllib/gonet/v3/codec"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	sys System //系统
)

type Hook func(msg *Message)

type System struct {
	sync.Once
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
	//服务端
	server Server
	//bee worker pool
	workers BeeWorkerPool
	//消息钩子
	mMsgHooks map[MessageID]Hook
}

func init() {
	log.SetPrefix("[gonet]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func init() {
	sys = System{
		SessionManager: SessionManager{
			pool: sync.Pool{
				New: func() interface{} {
					return reflect.New(sys.sessionType).Interface()
				},
			},
		},
		mMsgTypes: map[MessageID]reflect.Type{},
		mMsgIDs:   map[reflect.Type]MessageID{},
		mMsgHooks: map[MessageID]Hook{},
	}
}

func NewServer(opts ...options) Server {
	if sys.server == nil {
		panic(ErrorNoTransport)
	}
	option := Option{}
	for _, f := range opts {
		f(&option)
	}
	switch option.contentType {
	case Binary:
		sys.defaultCodec = codec.BinaryCodec{}
	case Xml:
		sys.defaultCodec = codec.XmlCodec{}
	case Protobuf:
		sys.defaultCodec = codec.ProtobufCodec{}
	default:
		sys.defaultCodec = codec.JsonCodec{}
	}
	cache := option.msgCache
	if cache == nil {
		cache = &MessageList{}
	}
	sys.workers = createBeeWorkerPool(option.workerPoolSize, cache)
	sys.server.(interface{ setAddr(string) }).setAddr(option.addr)
	return sys.server
}

//获取会话
func GetSession(id uint64) (Session, bool) {
	value, ok := sys.sessions.Load(id)
	if ok {
		return value.(Session), ok
	}
	return nil, false
}

//创建会话
func CreateSession() Session {
	obj := sys.pool.Get()
	//sys.incr = atomic.AddUint64(&sys.incr, 1)
	sys.store(atomic.AddUint64(&sys.incr, 1), obj)
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
	sys.del(session.ID())
	//回收
	sys.pool.Put(session)
}

//统计会话数量
func SessionCount() int {
	sum := 0
	sys.sessions.Range(func(key, value interface{}) bool {
		sum++
		return true
	})
	return sum
}

//广播会话
func Broadcast(msg interface{}) {
	sys.sessions.Range(func(_, item interface{}) bool {
		session, ok := item.(Session)
		if ok {
			session.Send(msg)
		}
		return true
	})
}

//映射消息体
func Route(msgID MessageID, msg interface{}, callback Hook) {
	sys.Lock()
	defer sys.Unlock()
	msgType := reflect.TypeOf(msg)
	if _, ok := sys.mMsgTypes[msgID]; ok {
		panic("error:Duplicate message id")
	}
	if msgType != nil {
		sys.mMsgIDs[msgType] = msgID
		sys.mMsgTypes[msgID] = msgType
	}
	if callback != nil {
		sys.mMsgHooks[msgID] = callback
	}
}

//获取消息ID
func GetMsgID(msg interface{}) (MessageID, bool) {
	msgID, ok := sys.mMsgIDs[reflect.TypeOf(msg)]
	return msgID, ok
}

//通消息id创建消息体
func CreateMsg(msgID MessageID) interface{} {
	if msg, ok := sys.mMsgTypes[msgID]; ok {
		return reflect.New(msg).Interface()
	}
	return nil
}

//初始化服务端
func RegisterServer(server Server, session interface{}) {
	sys.Once.Do(func() {
		sys.server = server
		sys.sessionType = reflect.TypeOf(session)
	})
}

//编码消息
func EncodeMessage(msg interface{}) ([]byte, error) {
	return sys.defaultCodec.Encode(msg)
}

// 解码消息
func DecodeMessage(msg interface{}, data []byte) error {
	return sys.defaultCodec.Decode(data, msg)
}

//缓存消息
func CacheMsg(msg *Message) {
	sys.workers.rcvMsgCh <- msg
}
