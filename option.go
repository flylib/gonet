package gonet

import (
	"github.com/zjllib/gonet/v3/transport"
	"time"
)

///////////////////////////////
/////    Option Func   ////////
//////////////////////////////

//options
type Option struct {
	//传输端
	transport transport.Transport
	//读写超时
	readDeadline, writeDeadline time.Duration
	//0意味着无限制
	maxSessionCount int
	//最小限制是1
	maxWorkerPoolSize int32
	//contentType support json/xml/binary/protobuf
	contentType string
	//worker pool size
	workerPoolSize int32
	//cache for messages
	msgCache MessageCache
}

type options func(o *Option)

func MaxSessions(max int) options {
	return func(o *Option) {
		o.maxSessionCount = max
	}
}

func MaxWorkerPoolSize(max int32) options {
	return func(o *Option) {
		o.maxWorkerPoolSize = max
	}
}

// Default content type of the client
func ContentType(ct string) options {
	return func(o *Option) {
		o.contentType = ct
	}
}

//传输协议
func WithTransport(t transport.Transport) options {
	return func(o *Option) {
		o.transport = t
	}
}

//cache for messages
func WithMessageCache(cache MessageCache) options {
	return func(o *Option) {
		o.msgCache = cache
	}
}
