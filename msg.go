package gonet

// 系统消息
const (
	MessageID_Invalid uint32 = iota
	MessageID_Connection_Connect
	MessageID_Connection_Close
)

var (
	msgNewConnection = &message{
		id: MessageID_Connection_Connect,
	}
)

type Event struct {
	Message IMessage
	Session ISession
}

type (
	EventHandler func(event Event)
)

type IMessage interface {
	ID() uint32
	Body() []byte
}

type message struct {
	id   uint32
	body []byte
}

func newErrorMessage(err error) *message {
	return &message{
		id:   MessageID_Invalid,
		body: []byte(err.Error()),
	}
}
func newCloseMessage(err error) *message {
	return &message{
		id:   MessageID_Connection_Close,
		body: []byte(err.Error()),
	}
}

func (m *message) ID() uint32 {
	return m.id
}

func (m *message) Body() []byte {
	return m.body
}
