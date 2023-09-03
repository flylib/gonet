package gonet

import (
	"container/list"
)

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

// 消息中间缓存层，为处理不过来的消息进行缓存
type IEventCache interface {
	Size() int
	Push(event IEvent)
	Pop() IEvent
}

// g默认的消息缓存队列
type MessageList struct {
	list.List
}

func (l *MessageList) Size() int {
	return l.List.Len()
}

func (l *MessageList) Push(msg IEvent) {
	l.List.PushFront(msg)
}

func (l *MessageList) Pop() IEvent {
	element := l.List.Back()
	if element == nil {
		return nil
	}
	l.List.Remove(element)
	return element.Value.(IEvent)
}
