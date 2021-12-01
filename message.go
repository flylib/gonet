package gonet

import (
	"reflect"
)

var (
	msgSessionConnect = SessionConnect{}
	msgSessionClose   = SessionClose{}
)

type MessageID uint32

//系统消息
const (
	MsgIDDecPoolSize uint32 = iota
	MsgIDSessionConnect
	MsgIDSessionClose
)

//消息体
type Message struct {
	Session
	ID   MessageID   `json:"id"`
	Body interface{} `json:"data"`
}

//会话
type SessionConnect struct {
}
type SessionClose struct {
}

//映射消息体
func RegisterMsg(msgID MessageID, msg interface{}) {
	sys.Lock()
	defer sys.Unlock()
	msgType := reflect.TypeOf(msg)
	if _, ok := sys.msgTypes[msgID]; ok {
		panic("error:Duplicate message id")
	} else {
		sys.msgTypes[msgID] = msgType
	}
	sys.msgIDs[msgType] = msgID
}

//获取消息ID
func GetMsgID(msg interface{}) (MessageID, bool) {
	msgID, ok := sys.msgIDs[reflect.TypeOf(msg)]
	return MessageID(msgID), ok
}

//通消息id创建消息体
func CreateMsg(msgID MessageID) interface{} {
	if msg, ok := sys.msgTypes[msgID]; ok {
		return reflect.New(msg).Interface()
	}
	return nil
}
