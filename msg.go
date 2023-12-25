package gonet

type IMessage interface {
	ID() uint32
	Body() []byte
	From() ISession
	UnmarshalTo(v any) error
}

type message struct {
	id      uint32
	body    []byte
	session ISession
}

func (m *message) ID() uint32 {
	return m.id
}

func (m *message) Body() []byte {
	return m.body
}

func (m *message) From() ISession {
	return m.session
}

func (m *message) UnmarshalTo(v any) error {
	return m.session.GetContext().Unmarshal(m.Body(), v)
}

type IEventHandler interface {
	OnConnect(ISession)
	OnClose(ISession, error)
	OnMessage(IMessage)
	OnError(ISession, error)
}
