package goNet

import (
	"math"
	"reflect"
	"time"
)

//消息接口
type Msg interface {
	//消息处理方法，每个消息要实现
	Handle(ses Session)
}

//心跳
type Ping struct {
	TimeStamp int64 `json:"time_stamp",xml:"time_stamp"`
}
type Pong struct {
	TimeStamp int64 `json:"time_stamp",xml:"time_stamp"`
}

func (p *Ping) Handle(session Session) {
	Log.Infof("session_%v ping at time=%v", session.ID(), time.Unix(p.TimeStamp, 0).String())
	session.Send(Pong{TimeStamp: time.Now().Unix()})
}
func (p *Pong) Handle(session Session) {
	Log.Infof("session_%v pong at time=%v", session.ID(), time.Unix(p.TimeStamp, 0).String())
}

//注册消息
func RegisterMsg(idx int32, msg interface{}) {
	if idx > math.MaxUint16 {
		Log.Panicf("msg.idx not allowed to be more than %v", math.MaxUint16)
	}
	more := idx - int32(len(msgTypes)) + 1
	if more > 0 {
		moreMsgTList := make([]reflect.Type, more)
		msgTypes = append(msgTypes, moreMsgTList...)
	}
	if msgTypes[idx] != nil {
		Log.Panicf("msg.idx %v is existed,please check out", idx)
	}
	t := reflect.TypeOf(msg)
	msgTypes[idx] = t
	msgIds[t] = idx
}

func GetMsgByIdx(idx int) Msg {
	return reflect.New(msgTypes[idx]).Interface().(Msg)
}

func GetMsgIdByType(t reflect.Type) int32 {
	return msgIds[t]
}

var (
	msgTypes = make([]reflect.Type, 8)      //index:msgID value:msgType
	msgIds   = make(map[reflect.Type]int32) //key:msgType value:msgID
)

//提交到协程池处理消息
func SubmitMsgToAntsPool(msg Msg, s Session) {
	if err := antsPool.Submit(func() { msg.Handle(s) }); err != nil {
		Log.Errorf("antsPool commit message error,reason is ", err.Error())
	}
}

func init() {
	RegisterMsg(0, Ping{})
	RegisterMsg(1, Pong{})
}
