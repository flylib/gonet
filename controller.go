package goNet

import (
	"github.com/astaxie/beego/logs"
	"time"
)

const SYSTEM_CONTROLLER_IDX = 0

var systemController Controller

func init() {
	systemController = &SystemController{}
}

//控制器
type Controller interface {
	OnMsg(session Session, msg interface{})
}

//系统控制模块
type SystemController struct {
}

func (*SystemController) OnMsg(session Session, msg interface{}) {
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

//替换系统控制器
func UpdateSystemController(c Controller) {
	if c == nil {
		return
	}
	systemController = c
}
