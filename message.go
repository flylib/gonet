package gonet

import (
	"container/list"
)

var _ MessageCache = new(MessageList)

type MessageID uint32

// 系统消息
const (
	SessionConnect MessageID = iota + 1
	SessionClose
)

type Head struct {
	session ISession
}

func (h *Head) setSession(session ISession) {
	h.session = session
}
func (h *Head) FromSession() ISession {
	return h.session
}

// 消息体
type Message struct {
	Head
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
	FromSession() ISession
}

// 消息中间缓存层，为处理不过来的消息进行缓存
type MessageCache interface {
	Size() int
	Push(msg *Message)
	Pop() *Message
}

// g默认的消息缓存队列
type MessageList struct {
	list.List
}

func (l *MessageList) Size() int {
	return l.List.Len()
}

func (l *MessageList) Push(msg *Message) {
	l.List.PushFront(msg)
}

func (l *MessageList) Pop() *Message {
	element := l.List.Back()
	if element == nil {
		return nil
	}
	l.List.Remove(element)
	return element.Value.(*Message)
}

type IPackageParser interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte) (IMessage, error)
}
