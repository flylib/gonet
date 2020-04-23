package goNet

import (
	"github.com/Quantumoffices/beego/logs"
	"time"
)

const SYSTEM_CONTROLLER_IDX = 0

var systemMsgController *SystemMsgController

func init() {
	systemMsgController = &SystemMsgController{}
}

//控制器
type MsgController interface {
	//消息处理接口
	ProcessMsg(session Session, msg interface{})
}

//系统控制模块
type SystemMsgController struct {
}

func (*SystemMsgController) ProcessMsg(session Session, msg interface{}) {
	switch data := msg.(type) {
	case *SessionConnect:
		logs.Info("session_%v connected", session.ID())
	case *SessionClose:
		logs.Warn("session_%v closed", session.ID())
	case *Ping:
		logs.Info("session_%v ping at time=%v", session.ID(), time.Unix(data.TimeStamp, 0).String())
		session.Send(Pong{TimeStamp: time.Now().Unix()})
	case *Pong:
		logs.Info("session_%v pong at time=%v", session.ID(), time.Unix(data.TimeStamp, 0).String())
	}
}
