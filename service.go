package gonet

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

var (
	peers = map[TransportProtocol]Service{}
)

type (
	//服务端
	Service interface {
		// 开启服务
		Start() error
		// 停止服务
		Stop() error
		//地址
		Addr() string
	}
	//端属性
	ServerIdentify struct {
		uuid string
		//地址
		addr string
	}
)

func (s *ServerIdentify) Addr() string {
	return s.addr
}

func (s *ServerIdentify) setAddr(addr string) {
	s.addr = addr
}

func RegisterPeer(peer Service) {
	peers[peer.(interface{ Type() PeerType }).Type()] = peer
}
