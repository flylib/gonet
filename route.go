package goNet

import (
	"github.com/astaxie/beego/logs"
	"time"
)

const System_Route_ID = 0

//系统消息路由
var sysRoute Route = &SystemController{}

//消息路由接口
type Route interface {
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

func UpdateSysRoute(c Route) {
	if c == nil {
		return
	}
	sysRoute = c
}
