package gonet

import "sync"

// IMessage represents a decoded network message.
type IMessage interface {
	ID() uint32
	Body() []byte
	From() ISession
	UnmarshalTo(v any) error
}

// IEventHandler handles connection lifecycle and message events.
type IEventHandler interface {
	OnConnect(ISession)
	OnClose(ISession, error)
	OnMessage(IMessage)
	OnError(ISession, error)
}

// msgPool reduces allocations for high-throughput message processing.
var msgPool = sync.Pool{
	New: func() any { return new(message) },
}

type message struct {
	id      uint32
	body    []byte
	session ISession
}

func newMessage(id uint32, body []byte, s ISession) *message {
	m := msgPool.Get().(*message)
	m.id = id
	m.body = body
	m.session = s
	return m
}

func releaseMessage(m *message) {
	m.id = 0
	m.body = nil
	m.session = nil
	msgPool.Put(m)
}

func (m *message) ID() uint32     { return m.id }
func (m *message) Body() []byte   { return m.body }
func (m *message) From() ISession { return m.session }

func (m *message) UnmarshalTo(v any) error {
	return m.session.GetContext().Unmarshal(m.body, v)
}
