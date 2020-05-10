package goNet

import (
	"github.com/astaxie/beego/logs"
	"github.com/panjf2000/ants/v2"
)

const (
	POOL_DEFAULT_SIZE = 1
)

var antsPool *ants.Pool

func init() {
	antsPool, _ = ants.NewPool(POOL_DEFAULT_SIZE)
}

//提交到协程池处理消息
func SubmitMsgToAntsPool(c Controller, s Session, msg interface{}) {
	if err := antsPool.Submit(func() {
		c.OnMsg(s, msg)
	}); err != nil {
		logs.Error("antsPool commit message error,reason is ", err.Error())
	}
}
