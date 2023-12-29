package gonet

type Message struct {
	id      uint32
	body    []byte
	session ISession
}

func (m *Message) ID() uint32 {
	return m.id
}

func (m *Message) Body() []byte {
	return m.body
}

func (m *Message) From() ISession {
	return m.session
}

func (m *Message) UnmarshalTo(v any) error {
	return defaultCtx.Unmarshal(m.Body(), v)
}
