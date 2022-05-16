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

type System struct {
	sync.Once
	ConnManager
	//message types
	msgTypes map[MessageID]reflect.Type
	//message ids
	msgIDs map[reflect.Type]MessageID
	//server types
	connType reflect.Type
	//消息编码器
	defaultCodec codec.Codec
	//服务端
	server Server
	//携程池
	workers WorkerPool
	//消息钩子
	mHandlers map[MessageID]SessionHandler
}

func init() {
	log.SetPrefix("[gonet]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func init() {
	sys = System{
		ConnManager: ConnManager{
			pool: sync.Pool{
				New: func() interface{} {
					return reflect.New(sys.connType).Interface()
				},
			},
		},
		msgTypes:  map[MessageID]reflect.Type{},
		msgIDs:    map[reflect.Type]MessageID{},
		mHandlers: map[MessageID]SessionHandler{},
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
	case codec.Binary:
		sys.defaultCodec = codec.BinaryCodec{}
	case codec.Xml:
		sys.defaultCodec = codec.XmlCodec{}
	case codec.Protobuf:
		sys.defaultCodec = codec.JsonCodec{}
	default:
		sys.defaultCodec = codec.JsonCodec{}
	}
	cache := option.msgCache
	if cache == nil {
		cache = &SessionCacheList{}
	}
	sys.workers = createWorkerPool(option.workerPoolSize, cache)
	sys.server.(interface{ setAddr(string) }).setAddr(option.addr)
	return sys.server
}

//获取会话
func GetConn(id uint64) (Connection, bool) {
	value, ok := sys.connections.Load(id)
	if ok {
		return value.(Connection), ok
	}
	return nil, false
}

//创建会话
func CreateConn() Connection {
	obj := sys.pool.Get()
	sys.store(atomic.AddUint64(&sys.incr, 1), obj)
	conn := obj.(Connection)
	return conn
}

//回收会话对象
func RecycleConn(conn Connection, err error) {
	CacheSession(&Session{
		Conn: conn,
		Msg:  msgConnClose,
	})
	//关闭
	conn.Close()
	//删除
	sys.del(conn.ID())
	//回收
	sys.pool.Put(conn)
}

//统计会话数量
func GetConnCount() int {
	sum := 0
	sys.connections.Range(func(key, value interface{}) bool {
		sum++
		return true
	})
	return sum
}

//广播会话
func Broadcast(msg interface{}) {
	sys.connections.Range(func(_, item interface{}) bool {
		conn, ok := item.(Connection)
		if ok {
			conn.Send(msg)
		}
		return true
	})
}

//映射消息体
func Route(msgID MessageID, msg interface{}, f SessionHandler) {
	sys.Lock()
	defer sys.Unlock()
	msgType := reflect.TypeOf(msg)
	if _, ok := sys.msgTypes[msgID]; ok {
		panic("error:Duplicate message id")
	}
	if msgType != nil {
		sys.msgIDs[msgType] = msgID
		sys.msgTypes[msgID] = msgType
	}
	if f != nil {
		sys.mHandlers[msgID] = f
	}
}

//获取消息ID
func GetMsgID(msg interface{}) (MessageID, bool) {
	msgID, ok := sys.msgIDs[reflect.TypeOf(msg)]
	return msgID, ok
}

//通消息id创建消息体
func CreateMsg(msgID MessageID) interface{} {
	if msg, ok := sys.msgTypes[msgID]; ok {
		return reflect.New(msg).Interface()
	}
	return nil
}

//初始化服务端
func RegisterServer(server Server, conn interface{}) {
	sys.Once.Do(func() {
		sys.server = server
		sys.connType = reflect.TypeOf(conn)
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
func CacheSession(s *Session) {
	sys.workers.sessionCh <- s
}

func GetCommonMsgNewConnMsg() *Message {
	return msgNewConn
}
