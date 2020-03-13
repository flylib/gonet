package goNet

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
	PEERTYPE_SERVER PeerType = "server" //服务端
	PEERTYPE_CLIENT PeerType = "client" //客户端
)

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

func NewPeer(opts ...Option) Peer {
	//parser options
	for _, opt := range opts {
		opt(Opts)
	}
	err := initAntsPool()
	if err != nil {
		panic(err)
	}
	p := peers[Opts.PeerType]
	p.(interface{ SetAddr(string) }).SetAddr(Opts.Addr)
	return p
}
