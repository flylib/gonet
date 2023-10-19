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
	_ ISessionIdentify = new(SessionIdentify)
	_ ISessionAbility  = new(SessionAbility)
	_ IPeerIdentify    = new(PeerIdentify)
)
