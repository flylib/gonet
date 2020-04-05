package goNet

import "time"

const SYSTEM_CONTROLLER_IDX = 0

var sysCtl Controller

func init() {
	sysCtl = &SystemController{}
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
		Log.Infof("session_%v connected", session.ID())
	case *SessionClose:
		Log.Warnf("session_%v closed", session.ID())
	case *Ping:
		Log.Infof("session_%v ping at time=%v", session.ID(), time.Unix(data.TimeStamp, 0).String())
		session.Send(Pong{TimeStamp: time.Now().Unix()})
	case *Pong:
		Log.Infof("session_%v pong at time=%v", session.ID(), time.Unix(data.TimeStamp, 0).String())
	}
}

//替换系统控制器
func UpdateSystemController(c Controller) {
	if c == nil {
		return
	}
	sysCtl = c
}
