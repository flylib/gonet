package goNet

import (
	"math"
	"reflect"
)

var (
	msgTypes            = make([]reflect.Type, 8)    //index:msgIdx value:msgType
	msgTypeIdx          = make(map[reflect.Type]int) //key:msgType value:msgIdx
	msgIdxControllerIdx = make(map[int]int)          //key:msgIdx value:controllerIdx
)

func init() {
	RegisterMsg(1, SYSTEM_CONTROLLER_IDX, SessionConnect{})
	RegisterMsg(2, SYSTEM_CONTROLLER_IDX, SessionClose{})
	RegisterMsg(3, SYSTEM_CONTROLLER_IDX, Ping{})
	RegisterMsg(4, SYSTEM_CONTROLLER_IDX, Pong{})
}

//心跳
type Ping struct {
	TimeStamp int64 `json:"time_stamp",xml:"time_stamp"`
}
type Pong struct {
	TimeStamp int64 `json:"time_stamp",xml:"time_stamp"`
}

//会话
type SessionConnect struct {
}
type SessionClose struct {
}

//注册消息
func RegisterMsg(msgIndex, controllerIndex int, msg interface{}) {
	if msgIndex > math.MaxUint16 {
		Log.Panicf("Message_%v index not allowed to be more than %v", msgIndex, math.MaxUint16)
	}
	more := msgIndex - len(msgTypes) + 1
	//extending
	if more > 0 {
		moreMsgTList := make([]reflect.Type, more)
		msgTypes = append(msgTypes, moreMsgTList...)
	}
	if msgTypes[msgIndex] != nil {
		Log.Panicf("message_%v duplicate registration error", msgIndex)
	}
	t := reflect.TypeOf(msg)
	msgTypes[msgIndex] = t
	msgTypeIdx[t] = msgIndex
	msgIdxControllerIdx[msgIndex] = controllerIndex
}

func GetMsgByIdx(msgIndex int) interface{} {
	return reflect.New(msgTypes[msgIndex]).Interface().(interface{})
}

func GetMsgIdxByType(t reflect.Type) int {
	return msgTypeIdx[t]
}

func GetMsgBelongToControllerIdx(msgIndex int) int {
	return msgIdxControllerIdx[msgIndex]
}
