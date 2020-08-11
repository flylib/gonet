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
	RegisterMsg(1, System_Route_ID, msgSessionConnect)
	RegisterMsg(2, System_Route_ID, msgSessionClose)
	RegisterMsg(3, System_Route_ID, msgPing)
	RegisterMsg(4, System_Route_ID, msgPong)
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

//@Param msg id
func FindMsg(msgID int) interface{} {
	return reflect.New(msgTypeList[msgID]).Interface()
}

//@Param msg type
func FindMsgID(t reflect.Type) int {
	return msgMap[t]
}

//@Param msg id
func FindRouteID(msgID int) int {
	return msgCtlMap[msgID]
}
