package gonet

// 系统消息
const (
	MessageID_Invalid uint32 = iota
	MessageID_Connection_Connect
	MessageID_Connection_Close
)

type (
	MessageHandler func(IMessage)
)

type IMessage interface {
	ID() uint32
	Body() []byte
	Session() ISession
	UnmarshalTo(v any) error
}

type message struct {
	id      uint32
	body    []byte
	session ISession
}

func newConnectionConnectMessage(s ISession) *message {
	return &message{
		id: MessageID_Connection_Connect,
	}
}

func newConnectionCloseMessage(s ISession, err error) *message {
	return &message{
		id: MessageID_Connection_Close,
	}
}

func newErrorMessage(s ISession, err error) *message {
	return &message{
		id:   MessageID_Invalid,
		body: []byte(err.Error()),
	}
}

func (m *message) ID() uint32 {
	return m.id
}

func (m *message) Body() []byte {
	return m.body
}

func (m *message) Session() ISession {
	return m.session
}

func (m *message) UnmarshalTo(v any) error {
	return m.session.Context().Unmarshal(m.Body(), v)
}
