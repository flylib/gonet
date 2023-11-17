package gonet

// 系统消息
const (
	MessageID_Invalid uint32 = iota
	MessageID_Connection_Connect
	MessageID_Connection_Close
	MessageID_Ping_Pong //ping pong
)

type (
	MessageHandler func(IMessage)
)

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

func newConnectionConnectMessage(s ISession) *message {
	return &message{
		id:      MessageID_Connection_Connect,
		session: s,
	}
}

func newConnectionCloseMessage(s ISession, err error) *message {
	return &message{
		id:      MessageID_Connection_Close,
		session: s,
	}
}

func newErrorMessage(s ISession, err error) *message {
	return &message{
		id:      MessageID_Invalid,
		body:    []byte(err.Error()),
		session: s,
	}
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
