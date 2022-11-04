package gonet

import "container/list"

var _ MessageCache = new(MessageList)

type MessageID uint32

//系统消息
const (
	SessionConnect MessageID = iota + 1
	SessionClose
)

//消息体
type Message struct {
	Session
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

const (
	Json     = "json"
	Binary   = "binary"
	Protobuf = "protobuf"
	Xml      = "xml"
)

type Codec interface {
	Encode(v interface{}) (data []byte, err error)
	Decode(data []byte, vObj interface{}) error
	Type() string
}
