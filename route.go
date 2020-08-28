package goNet

import (
	"github.com/astaxie/beego/logs"
	"time"
)

//默认ID
const DefaultRouteID = 0

//系统消息路由
var defaultRoute Route = &DefaultRoute{}

//消息路由接口
type Route interface {
	OnMsg(session Session, msg interface{})
}

//系统控制模块
type DefaultRoute struct {
}

func (*DefaultRoute) OnMsg(session Session, msg interface{}) {
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
	defaultRoute = c
}

//获取消息所在route
func FindMsgOnRoute(msgID int) int {
	return mMsgRoute[msgID]
}
