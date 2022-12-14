package gonet

import (
	"container/list"
	"github.com/zjllib/gonet/v3/transport"
)

var _ MessageCache = new(MessageList)

type MessageID uint32

//系统消息
const (
	SessionConnect MessageID = iota + 1
	SessionClose
)

type Head struct {
	session transport.ISession
}

func (h *Head) setSession(session transport.ISession) {
	h.session = session
}
func (h *Head) GetSession() transport.ISession {
	return h.session
}

//消息体
type Message struct {
	Head
	ID   MessageID   `json:"id"`
	Body interface{} `json:"data"`
}

//消息中间缓存层，为处理不过来的消息进行缓存
type MessageCache interface {
	Size() int
	Push(msg *Message)
	Pop() *Message
}

//g默认的消息缓存队列
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
