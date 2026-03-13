package gonet

import (
	"github.com/flylib/interface/codec"
	ilog "github.com/flylib/interface/log"
)

// config holds non-generic configuration for Context.
type config struct {
	eventHandler    IEventHandler
	maxSessionCount int
	poolCfg         poolConfig
	codec.ICodec
	ilog.ILogger
	INetPackager
}

// AppContext[S] is the core framework context, generic over the session type S.
//
// Usage:
//
//	ctx := gonet.NewAppContext(
//	    func() *tcp.Session { return new(tcp.Session) },
//	    gonet.WithEventHandler(handler),
//	    gonet.MustWithCodec(codec),
//	    gonet.MustWithLogger(logger),
//	)
type AppContext[S SessionConstraint] struct {
	config
	sessions *sessionManager[S]
	routines *GoroutinePool
}

// NewAppContext creates a new Context. factory must return a new, zero-value session.
func NewAppAppContext[S SessionConstraint](factory func() S, options ...Option) *AppContext[S] {
	cfg := config{INetPackager: &DefaultNetPackager{}}
	for _, f := range options {
		f(&cfg)
	}
	if cfg.ICodec == nil {
		panic("gonet: ICodec is required, use MustWithCodec()")
	}
	if cfg.ILogger == nil {
		panic("gonet: ILogger is required, use MustWithLogger()")
	}
	if cfg.eventHandler == nil {
		panic("gonet: IEventHandler is required, use WithEventHandler()")
	}
	ctx := &AppContext[S]{config: cfg}
	ctx.sessions = newSessionManager(factory)
	ctx.routines = newGoroutinePool(cfg.poolCfg, cfg.ILogger, cfg.eventHandler)
	return ctx
}

// --- IContext implementation ---

func (c *AppContext[S]) GetEventHandler() IEventHandler { return c.eventHandler }

func (c *AppContext[S]) PushGlobalMessageQueue(msg IMessage) {
	if c.poolCfg.queueSize == 0 {
		// No pool: session's own goroutine handles the message synchronously.
		c.eventHandler.OnMessage(msg)
		if m, ok := msg.(*message); ok {
			recycleMessage(m)
		}
		return
	}
	c.routines.push(msg)
}

func (c *AppContext[S]) GetSession(id uint64) (ISession, bool) {
	return c.sessions.getAlive(id)
}

func (c *AppContext[S]) SessionCount() int32 {
	return c.sessions.count()
}

// RecycleSession closes the session, clears it, and returns it to the idle pool.
// Bug fix: ID is saved before Clear() resets it to 0.
func (c *AppContext[S]) RecycleSession(session ISession) {
	id := session.ID() // save before Clear resets it
	_ = session.Close()
	s := session.(S)
	s.Clear()
	c.sessions.removeAlive(id)
	c.sessions.putIdle(s)
}

// UpdateSessionID atomically moves a session to a new ID in the alive map.
func (c *AppContext[S]) UpdateSessionID(session ISession, newID uint64) {
	oldID := session.ID()
	c.sessions.alive.Delete(oldID)
	c.sessions.alive.Store(newID, session)
}

func (c *AppContext[S]) Broadcast(msgID uint32, msg any) {
	c.sessions.alive.Range(func(_, v any) bool {
		if s, ok := v.(ISession); ok {
			_ = s.Send(msgID, msg)
		}
		return true
	})
}

func (c *AppContext[S]) Package(s ISession, msgID uint32, v any) ([]byte, error) {
	return c.INetPackager.Package(s, msgID, v)
}

func (c *AppContext[S]) UnPackage(s ISession, data []byte) (IMessage, int, error) {
	return c.INetPackager.UnPackage(s, data)
}

func (c *AppContext[S]) Marshal(v any) ([]byte, error) {
	return c.ICodec.Marshal(v)
}

func (c *AppContext[S]) Unmarshal(data []byte, v any) error {
	return c.ICodec.Unmarshal(data, v)
}

// GetIdleSession retrieves a session from the pool and registers it as alive.
// Returns (session, false) if maxSessionCount has been reached.
func (c *AppContext[S]) GetIdleSession() (S, bool) {
	if c.maxSessionCount > 0 && int(c.sessions.count()) >= c.maxSessionCount {
		var zero S
		return zero, false
	}
	s := c.sessions.getIdle()
	s.Clear()
	s.WithContext(c)
	c.sessions.addAlive(s)
	return s, true
}

// GetLogger returns the configured logger.
func (c *AppContext[S]) GetLogger() ilog.ILogger {
	return c.ILogger
}
