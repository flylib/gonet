package goNet

import "github.com/astaxie/beego/logs"

var (
	peers = map[PeerType]Peer{}
)

type (
	//端
	Peer interface {
		// 开启服务
		Start()

		// 停止服务
		Stop()
	}
	//端属性
	PeerIdentify struct {
		//地址
		addr string
		//类型
		peerType PeerType
	}
	//端类型
	PeerType string
)

const (
	PeertypeServer PeerType = "server" //服务端
	PeertypeClient PeerType = "client" //客户端
)

func init() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
}

func (p *PeerIdentify) Addr() string {
	return p.addr
}
func (p *PeerIdentify) SetAddr(addr string) {
	p.addr = addr
}
func (p *PeerIdentify) Type() PeerType {
	return p.peerType
}
func (p *PeerIdentify) SetType(t PeerType) {
	p.peerType = t
}

func RegisterPeer(peer Peer) {
	peers[peer.(interface{ Type() PeerType }).Type()] = peer
}

func NewServer(addr string, opts ...options) Peer {
	peer := peers[PeertypeServer]
	peer.(interface{ SetAddr(string) }).SetAddr(addr)
	option := Option{}
	for _, f := range opts {
		f(&option)
	}
	initWorkerPool(option)
	return peer
}

func NewClient(addr string, opts Option) Peer {
	peer := peers[PeertypeClient]
	peer.(interface{ SetAddr(string) }).SetAddr(addr)
	peer.(interface{ SetOptions(Option) }).SetOptions(opts)
	return peer
}
