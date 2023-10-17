package gonet

import "container/list"

var (
	_ IMessageCache = new(DefaultMessageCacheList)
)

type (
	MessageHandler func(ISession, IMessage)
	MessageID      uint32
)

// 系统消息
const (
	MessageID_Invalid MessageID = iota
	MessageID_SessionConnect
	MessageID_SessionClose
)

type Message struct {
	id   MessageID
	body []byte
}

func newInvalidMessage() *Message {
	return &Message{
		id: MessageID_Invalid,
	}
}

func newSessionConnectMessage(s ISession) *Message {
	return &Message{
		id: MessageID_SessionConnect,
	}
}

func newSessionCloseMessage(s ISession, err error) *Message {
	return &Message{
		id: MessageID_SessionConnect,
	}
}

func (m *Message) ID() MessageID {
	return m.id
}

func (m *Message) Payload() []byte {
	return m.body
}

// g默认的消息缓存队列
type DefaultMessageCacheList struct {
	list.List
}

func (l *DefaultMessageCacheList) Size() int {
	return l.List.Len()
}

func (l *DefaultMessageCacheList) Push(msg IMessage) {
	l.List.PushFront(msg)
}

func (l *DefaultMessageCacheList) Pop() IMessage {
	element := l.List.Back()
	if element == nil {
		return nil
	}
	l.List.Remove(element)
	return element.Value.(IMessage)
}
