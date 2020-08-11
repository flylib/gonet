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
		//配置
		options Options
	}
	//端类型
	PeerType string
)

const (
	PEERTYPE_SERVER PeerType = "server" //服务端
	PEERTYPE_CLIENT PeerType = "client" //客户端
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
func (p *PeerIdentify) SetOptions(o Options) {
	p.options = o
}
func (p *PeerIdentify) Options(o Options) {
	p.options = o
}

func RegisterPeer(peer Peer) {
	peers[peer.(interface{ Type() PeerType }).Type()] = peer
}

func NewPeer(opts Options) Peer {
	peer, ok := peers[opts.PeerType]
	if !ok {
		panic(opts.PeerType + "does not exist")
	}
	peer.(interface{ SetAddr(string) }).SetAddr(opts.Addr)
	peer.(interface{ SetOptions(Options) }).SetOptions(opts)
	//init worker pool
	InitWorkerPool(opts.PanicHandler, opts.EventChanSize)
	return peer
}
