package goNet

import (
	"math"
	"reflect"
)

var (
	arrMsgTypes = make([]reflect.Type, 8)    //index:msgIdx value:msgType
	maxMsgIndex = 0                          //最大消息索引
	mMsgType    = make(map[reflect.Type]int) //key:msgType value:msgID
	mMsgRoute   = make(map[int]int)          //key:msgID value:controllerID
)

var (
	msgSessionConnect = SessionConnect{}
	msgSessionClose   = SessionClose{}
	msgPing           = Ping{}
	msgPong           = Pong{}
)

func init() {
	RegisterMsg(1, DefaultRouteID, msgSessionConnect)
	RegisterMsg(2, DefaultRouteID, msgSessionClose)
	RegisterMsg(3, DefaultRouteID, msgPing)
	RegisterMsg(4, DefaultRouteID, msgPong)
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
	between := msgID - len(arrMsgTypes) + 1
	//扩容
	if between > 0 {
		more := make([]reflect.Type, between)
		arrMsgTypes = append(arrMsgTypes, more...)
	}
	if arrMsgTypes[msgID] != nil {
		panic("Duplicate message")
	}
	t := reflect.TypeOf(msg)
	arrMsgTypes[msgID] = t
	mMsgType[t] = msgID
	mMsgRoute[msgID] = controllerIndex
	maxMsgIndex = len(arrMsgTypes) - 1
}

//实例化消息
func InstanceMsg(msgID int) (interface{}, error) {
	if msgID > maxMsgIndex || arrMsgTypes[msgID] == nil {
		return nil, ErrNotFoundMsg
	}
	return reflect.New(arrMsgTypes[msgID]).Interface(), nil
}

//获取消息ID
func GetMsgID(t reflect.Type) int {
	return mMsgType[t]
}

////////////////////
////   EVENT   ////
///////////////////

//事件分类
const (
	EventNetWorkIO  EventType = iota //default,网络i/o
	EventWorkerExit                  //退出worker
	EventWorkerAdd                   //新增worker
)

type EventType int8

//事件
type Event struct {
	eventType EventType   //事件分类
	from      Session     //来自
	route     Route       //路由(处理器)
	data      interface{} //消息-事件内容
}

//创建事件
func CreateEvent(t EventType, session Session, route Route, data interface{}) Event {
	return Event{
		eventType: t,
		from:      session,
		route:     route,
		data:      data,
	}
}
