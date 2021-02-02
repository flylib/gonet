package goNet

import (
	"reflect"
)

var (
	//msgID:msgType
	mMsg = map[uint32]reflect.Type{}
	//msgType:msgID
	mMsgType = map[reflect.Type]uint32{}
	//msgID:sceneID
	mScene = map[uint32]uint8{}
)

var (
	msgSessionConnect = SessionConnect{}
	msgSessionClose   = SessionClose{}
	msgPing           = Ping{}
	msgPong           = Pong{}
)

func init() {
	RegisterMsg(1, DefaultSceneID, msgSessionConnect)
	RegisterMsg(2, DefaultSceneID, msgSessionClose)
	RegisterMsg(3, DefaultSceneID, msgPing)
	RegisterMsg(4, DefaultSceneID, msgPong)
}

//消息体
type Msg struct {
	Session
	SceneID uint8       `json:"scene_id"` //对应场景
	ID      uint32      `json:"id"`
	Data    interface{} `json:"data"`
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

//绑定场景消息
func RegisterMsg(sceneID uint8, msgID uint32, msg interface{}) {
	mScene[msgID] = sceneID
	msgType := reflect.TypeOf(msg)
	mMsg[msgID] = msgType
	mMsgType[msgType] = msgID
}

//获取消息ID
func GetMsgID(msg interface{}) uint32 {
	return mMsgType[reflect.TypeOf(msg)]
}

//获取消息所在场景ID
func GetMsgSceneID(msgID uint32) uint8 {
	return mScene[msgID]
}
func GetMsg(msgID uint32) interface{} {
	return reflect.New(mMsg[msgID]).Interface()
}

//事件分类
const (
	EventNetWorkIO  EventType = iota //default,网络io
	EventWorkerExit                  //退出worker
	EventWorkerAdd                   //新增worker
)

type EventType int8

//事件
type Event struct {
	eventType EventType //事件分类
	Actor     Actor     //路由(处理器)
	context
}
