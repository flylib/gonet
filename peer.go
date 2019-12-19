package goNet

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
		_type string
	}
)

func (p *PeerIdentify) Addr() string {
	return p.addr
}

func (p *PeerIdentify) SetAddr(addr string) {
	p.addr = addr
}

func (p *PeerIdentify) Type() string {
	return p._type
}

func (p *PeerIdentify) SetType(t string) {
	p._type = t
}

var (
	peers =map[string]Peer{}
)

func RegisterPeer(peer Peer) {
	peers[peer.(interface{Type() string}).Type()]=peer
}

func NewPeer(peertype,addr string) Peer {
	p:=peers[peertype]
	p.(interface {SetAddr(string)}).SetAddr(addr)
	return p
}
