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
	//SERVER
	server transport.IServer
	//CLIENT
	client transport.IClient
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
	//service name
	serviceName string
}

type options func(o *Option)

//server
func Server(s transport.IServer) options {
	return func(o *Option) {
		o.server = s
	}
}

//client
func Client(c transport.IClient) options {
	return func(o *Option) {
		o.client = c
	}
}

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

//cache for messages
func WithMessageCache(cache MessageCache) options {
	return func(o *Option) {
		o.msgCache = cache
	}
}

//cache for messages
func ServiceName(name string) options {
	return func(o *Option) {
		o.serviceName = name
	}
}
