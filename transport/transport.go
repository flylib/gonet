package transport

import "reflect"

type TransportProtocol string

const (
	TCP  TransportProtocol = "tcp"
	KCP  TransportProtocol = "kcp"
	UDP  TransportProtocol = "udp"
	WS   TransportProtocol = "websocket"
	HTTP TransportProtocol = "http"
	QUIC TransportProtocol = "quic"
	RPC  TransportProtocol = "rpc"
)

type (

	//传输协议
	Transport interface {
		// 启动监听
		Listen() error
		// 停止服务
		Stop() error
		// 地址
		Addr() string
		// 会话类型
		SessionType() reflect.Type
	}
	//端属性
	TransportIdentify struct {
		uuid string
		//地址
		addr string
	}
)

func (s *TransportIdentify) Addr() string {
	return s.addr
}

func (s *TransportIdentify) SetAddr(addr string) {
	s.addr = addr
}
