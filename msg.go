package gonet

import (
	"reflect"
)

var (
	//msgID:msgType
	mMsg = map[uint32]reflect.Type{}
	//msgType:msgID
	mMsgIDs = map[reflect.Type]uint32{}
	//msgID:sceneID
	mSceneIDs = map[uint32]uint8{}
)

var (
	msgSessionConnect = SessionConnect{}
	msgSessionClose   = SessionClose{}
)

//系统消息
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
	if msgID < 1 {
		panic("must be msgID>0")
	}
	if _, ok := mSceneIDs[msgID]; ok {
		panic("msg duplicate")
	}
	mSceneIDs[msgID] = sceneID
	msgType := reflect.TypeOf(msg)
	mMsg[msgID] = msgType
	mMsgIDs[msgType] = msgID
}

//获取消息ID
func GetMsgID(msg interface{}) uint32 {
	return mMsgIDs[reflect.TypeOf(msg)]
}

//获取消息所在场景ID
func GetMsgSceneID(msgID uint32) uint8 {
	return mSceneIDs[msgID]
}
func GetMsg(msgID uint32) interface{} {
	if msg, ok := mMsg[msgID]; ok {
		return reflect.New(msg).Interface()
	}
	return nil
}
