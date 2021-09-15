package gonet

import (
	"time"
)

type ProtoCol string

const (
	TCP  ProtoCol = "tcp"
	KCP  ProtoCol = "kcp"
	UDP  ProtoCol = "udp"
	WS   ProtoCol = "websocket"
	HTTP ProtoCol = "http"
	QUIC ProtoCol = "quic"
	RPC  ProtoCol = "rpc"
)

//options
type Option struct {
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
