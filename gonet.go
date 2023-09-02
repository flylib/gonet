package gonet

import (
	"github.com/zjllib/gonet/v3/codec"
	"github.com/zjllib/gonet/v3/codec/binary"
	"github.com/zjllib/gonet/v3/codec/json"
	"github.com/zjllib/gonet/v3/codec/protobuf"
	"github.com/zjllib/gonet/v3/codec/xml"
	"log"
	"reflect"
)

//var (
//	goNetContext Context //上下文
//)

func init() {
	log.SetPrefix("[gonet]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func init() {

}

type Hook func(msg *Message)

func NewContext(opts ...options) *Context {
	option := Option{}
	for _, f := range opts {
		f(&option)
	}
	c := &Context{
		sessionMgr: newSessionManager(option.server.SessionType()),
		mMsgTypes:  map[MessageID]reflect.Type{},
		mMsgIDs:    map[reflect.Type]MessageID{},
		mMsgHooks:  map[MessageID]Hook{},
	}
	//传输协议
	c.server = option.server
	c.sessionType = option.server.SessionType()
	if option.serviceName == "" {
		option.serviceName = ""
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
	c.workers = createBeeWorkerPool(option.workerPoolSize, cache)
	return c
}
