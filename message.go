package goNet

import (
	"math"
	"reflect"
)

var (
	msgTypeList = make([]reflect.Type, 8)    //index:msgIdx value:msgType
	msgMap      = make(map[reflect.Type]int) //key:msgType value:msgIdx
	msgCtlMap   = make(map[int]int)          //key:msgIdx value:controllerIdx
)
var (
	msgSessionConnect = SessionConnect{}
	msgSessionClose   = SessionClose{}
	msgPing           = Ping{}
	msgPong           = Pong{}
)

func init() {
	RegisterMsg(1, SYSTEM_CONTROLLER_IDX, msgSessionConnect)
	RegisterMsg(2, SYSTEM_CONTROLLER_IDX, msgSessionClose)
	RegisterMsg(3, SYSTEM_CONTROLLER_IDX, msgPing)
	RegisterMsg(4, SYSTEM_CONTROLLER_IDX, msgPong)
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
func RegisterMsg(msgID, controllerIndex int, msg interface{}) {
	if msgID > math.MaxUint16 {
		panic("msg index over allowed range")
	}
	more := msgID - len(msgTypeList) + 1
	//extending
	if more > 0 {
		moreMsgTList := make([]reflect.Type, more)
		msgTypeList = append(msgTypeList, moreMsgTList...)
	}
	if msgTypeList[msgID] != nil {
		panic("msg existed!")
	}
	t := reflect.TypeOf(msg)
	msgTypeList[msgID] = t
	msgMap[t] = msgID
	msgCtlMap[msgID] = controllerIndex
}

func FindMsg(msgID int) interface{} {
	return reflect.New(msgTypeList[msgID]).Interface()
}

func FindMsgIDByType(t reflect.Type) int {
	return msgMap[t]
}

func GetMsgBelongToControllerIdx(msgIndex int) int {
	return msgCtlMap[msgIndex]
}
