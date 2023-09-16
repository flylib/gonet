package gonet

var (
	_                 IEventCache = new(MessageList)
	NewSessionMessage             = Message{
		id: SessionConnect,
	}
	SessionCloseMessage = Message{
		id: SessionClose,
	}
)

type MessageID uint32

// 系统消息
const (
	SessionConnect MessageID = iota + 1
	SessionClose
)

// 消息体
type Message struct {
	id      MessageID
	body    any
	rawData []byte
}

func (m Message) ID() MessageID {
	return m.id
}

func (m Message) Body() any {
	return m.body
}

func (m Message) RawData() []byte {
	return m.rawData
}

type IMessage interface {
	ID() MessageID
	Body() any
	RawData() []byte
}
