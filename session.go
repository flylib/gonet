package gonet

import (
	"net"
	"sync"
	"sync/atomic"
)

// IContext is the interface through which sessions interact with the framework.
// *AppContext[S] implements this interface.
type IContext interface {
	GetEventHandler() IEventHandler
	PushGlobalMessageQueue(msg IMessage)
	RecycleSession(session ISession)
	UpdateSessionID(session ISession, newID uint64)
	GetSession(id uint64) (ISession, bool)
	SessionCount() int32
	Broadcast(msgID uint32, msg any)
	Package(s ISession, msgID uint32, v any) ([]byte, error)
	UnPackage(s ISession, data []byte) (IMessage, int, error)
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	NewMsg(id uint32, body []byte, s ISession) IMessage
	RecycleMsg(msg IMessage)
}

// BodyBufferProvider provides pooled buffers for message bodies.
type BodyBufferProvider interface {
	GetBodyBuffer(size int) []byte
	RecycleBodyBuffer(buf []byte)
}

// ISession is the public session interface exposed to users.
type ISession interface {
	ID() uint64
	Close() error
	Send(msgID uint32, msg any) error
	RemoteAddr() net.Addr
	Store(value any)
	Load() any
	GetContext() IContext
}

// SessionConstraint is the type constraint for AppContext[S].
// Embed SessionCommon in your session struct to satisfy this constraint.
type SessionConstraint interface {
	ISession
	SetID(id uint64)
	Clear()
	WithContext(c IContext)
}

type invalidData struct{}

var zeroData = invalidData{}

// SessionCommon provides base session functionality.
// Embed this in your transport-specific session struct.
type SessionCommon struct {
	ctx IContext
	atomic.Value
	sync.Mutex
	id uint64
}

func (s *SessionCommon) ID() uint64             { return s.id }
func (s *SessionCommon) GetContext() IContext   { return s.ctx }
func (s *SessionCommon) WithContext(c IContext) { s.ctx = c }
func (s *SessionCommon) SetID(id uint64)        { s.id = id }

func (s *SessionCommon) Clear() {
	s.ctx = nil
	s.id = 0
	s.Value.Store(&zeroData)
}

// UpdateID moves the session to a new ID in the alive map.
// Pass the outer ISession (the concrete type embedding SessionCommon).
func (s *SessionCommon) UpdateID(outer ISession, newID uint64) {
	if s.ctx != nil {
		s.ctx.UpdateSessionID(outer, newID)
	}
	s.id = newID
}
