package goNet

import (
	"github.com/panjf2000/ants"
	"github.com/sirupsen/logrus"
	"math"
	"reflect"
	"time"
)

//消息接口
type Message interface {
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
	logrus.Infof("session_%v ping at time=%v", session.ID(), time.Unix(p.TimeStamp, 0).String())
	session.Send(Pong{
		TimeStamp: time.Now().Unix(),
	})
}
func (p *Pong) Handle(session Session) {
	logrus.Infof("session_%v pong at time=%v", session.ID(), time.Unix(p.TimeStamp, 0).String())
}

//注册消息
func RegisterMessage(idx int32, msg interface{}) {
	if idx > math.MaxUint16 {
		logrus.Panicf("msg.idx not allowed to be more than %v", math.MaxUint16)
	}

	more := idx - int32(len(msgTList)) + 1
	if more > 0 {
		moreMsgTList := make([]reflect.Type, more)
		msgTList = append(msgTList, moreMsgTList...)
	}

	if msgTList[idx] != nil {
		logrus.Panicf("msg.idx %v is existed,please check out", idx)
	}

	t := reflect.TypeOf(msg)
	//if _, ok := reflect.New(t).Interface().(Message); !ok {
	//	logrus.Panic(t.Name() + ":not a Message type")
	//}
	msgTList[idx] = t
	msgTID[t] = idx
}

func GetMessageByIdx(idx int) Message {
	return reflect.New(msgTList[idx]).Interface().(Message)
}

func GetMessageID(t reflect.Type) int32 {
	return msgTID[t]
}

var (
	msgTList = make([]reflect.Type, 8)
	msgTID   = make(map[reflect.Type]int32)

	antsPool, _ = ants.NewPool(10)
)

//重置协程池大小
func ResetAnstsPoolSize(size int) {
	antsPool.Tune(size)
}

//提交到协程池处理消息
func HandleMessage(msg Message, s Session) {
	err := antsPool.Submit(func() {
		msg.Handle(s)
	})
	if err != nil {
		logrus.Errorf("antsPool commit message error,reason is ", err.Error())
	}
}

func init() {
	//0
	RegisterMessage(0, Ping{})
	//1
	RegisterMessage(1, Pong{})
}
