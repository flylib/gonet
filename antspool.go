package goNet

import (
	"github.com/astaxie/beego/logs"
	"github.com/panjf2000/ants/v2"
)

const (
	Default_Pool_Size = 10
)

var AntsPool *ants.Pool

func init() {
	AntsPool, _ = ants.NewPool(Default_Pool_Size)
}

//提交到协程池处理消息
func SubmitMsgToAntsPool(c Controller, s Session, msg interface{}) {
	if err := AntsPool.Submit(func() {
		c.OnMsg(s, msg)
	}); err != nil {
		logs.Error("antsPool commit message error,reason is ", err.Error())
	}
}
