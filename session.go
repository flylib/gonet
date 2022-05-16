package gonet

import "container/list"

//单次会话
type Session struct {
	Connection          //来自链接
	Msg        *Message //消息
}

//单次会话处理钩子
type SessionHandler func(msg *Session)

//会话中间缓存层，为处理不过来的会话进行缓存
type SessionCache interface {
	Size() int
	Push(msg *Session)
	Pop() *Session
}

//默认的消息缓存队列
type MessageList struct {
	list.List
}

func (l *MessageList) Size() int {
	return l.List.Len()
}

func (l *MessageList) Push(msg *Session) {
	l.List.PushFront(msg)
}

func (l *MessageList) Pop() *Session {
	element := l.List.Back()
	if element == nil {
		return nil
	}
	l.List.Remove(element)
	return element.Value.(*Session)
}
