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
	SessionManager
	//message types
	msgTypes map[MessageID]reflect.Type
	//message ids
	msgIDs map[reflect.Type]MessageID
	//server types
	sessionType reflect.Type
	//消息编码器
	defaultCodec codec.Codec
	//服务端
	server Server
}

func init() {
	log.SetPrefix("[gonet]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
	sys = System{
		SessionManager: SessionManager{
			pool: sync.Pool{
				New: func() interface{} {
					return reflect.New(sys.sessionType).Interface()
				},
			},
		},
	}
}

func NewServer(opts ...options) Server {
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
	if sys.server == nil {
		panic(ErrorNoTransport)
	}
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
	sys.store(atomic.AddUint64(&sys.incr, 1), obj)
	session := obj.(Session)
	return session
}

//回收会话对象
func RecycleSession(session Session) {
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
func RegisterMsg(msgID MessageID, msg interface{}) {
	sys.Lock()
	defer sys.Unlock()
	msgType := reflect.TypeOf(msg)
	if _, ok := sys.msgTypes[msgID]; ok {
		panic("error:Duplicate message id")
	} else {
		sys.msgTypes[msgID] = msgType
	}
	sys.msgIDs[msgType] = msgID
}

//获取消息ID
func GetMsgID(msg interface{}) (MessageID, bool) {
	msgID, ok := sys.msgIDs[reflect.TypeOf(msg)]
	return MessageID(msgID), ok
}

//通消息id创建消息体
func CreateMsg(msgID MessageID) interface{} {
	if msg, ok := sys.msgTypes[msgID]; ok {
		return reflect.New(msg).Interface()
	}
	return nil
}

//初始化服务端
func RegisterServer(server Server, session Session) {
	sys.Do(func() {
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
