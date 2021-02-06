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
)

const (
	MsgIDDecPoolSize uint32 = iota
	MsgIDSessionConnect
	MsgIDSessionClose
)

//消息体
type Msg struct {
	Session
	SceneID uint8       `json:"scene_id"` //对应场景
	ID      uint32      `json:"id"`
	Data    interface{} `json:"data"`
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
