package gonet

import "container/list"

var (
	_ IMessageCache = new(DefaultMessageCacheList)
)

type (
	MessageHandler func(IMessage)
	MessageID      uint32
)

// 系统消息
const (
	MessageID_SessionConnect MessageID = iota + 1
	MessageID_SessionClose
)

type IMessage interface {
	ID() MessageID
	Body() any
	RawData() []byte
	FromSession() ISession
}

// 消息体
type Message struct {
	id      MessageID
	body    any
	rawData []byte
	session ISession
}

func newSessionConnectMessage(s ISession) *Message {
	return &Message{
		id:      MessageID_SessionConnect,
		session: s,
	}
}

func newSessionCloseMessage(s ISession, err error) *Message {
	return &Message{
		id:      MessageID_SessionConnect,
		body:    err,
		session: s,
	}
}

func (m *Message) ID() MessageID {
	return m.id
}

func (m *Message) Body() any {
	return m.body
}

func (m *Message) RawData() []byte {
	return m.rawData
}
func (m *Message) FromSession() ISession {
	return m.session
}

// 消息中间缓存层，为处理不过来的消息进行缓存
type IMessageCache interface {
	Size() int
	Push(IMessage)
	Pop() IMessage
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
