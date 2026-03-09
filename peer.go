package gonet

// Transport layer interfaces.
type (
	// IServer is the server-side transport interface.
	IServer interface {
		Listen(addr string) error
		Close() error
		Addr() string
	}
	// IClient is the client-side transport interface.
	IClient interface {
		Dial(addr string) (ISession, error)
		Close() error
	}
)

// PeerCommon[S] provides shared server/client fields.
// Transport implementations embed this with their concrete session type.
type PeerCommon[S SessionConstraint] struct {
	ctx  *Context[S]
	addr string
}

func (p *PeerCommon[S]) Addr() string              { return p.addr }
func (p *PeerCommon[S]) SetAddr(addr string)       { p.addr = addr }
func (p *PeerCommon[S]) WithContext(c *Context[S]) { p.ctx = c }
func (p *PeerCommon[S]) GetCtx() *Context[S]       { return p.ctx }
