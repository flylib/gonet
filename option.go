package gonet

import (
	"time"
)

///////////////////////////////
/////    Option Func   ////////
//////////////////////////////

//options
type Option struct {
	//传输协议
	tpl TransportProtocol
	//关联地址
	addr string
	//读写超时
	readDeadline, writeDeadline time.Duration
	//0意味着无限制
	maxSessionCount int
	//最小限制是1
	maxWorkerPoolSize int32
}

type options func(o *Option)

func WithMaxSessions(max int) options {
	return func(o *Option) {
		o.maxSessionCount = max
	}
}

func WithMaxWorkerPoolSize(max int32) options {
	return func(o *Option) {
		o.maxWorkerPoolSize = max
	}
}
