package goNet

import (
	"github.com/panjf2000/ants/v2"
)

const (
	POOL_DEFAULT_SIZE = 1
)

var antsPool *ants.Pool

func initAntsPool() error {
	var err error
	if Opts.PoolSize <= 0 {
		antsPool, err = ants.NewPool(POOL_DEFAULT_SIZE)
	} else {
		antsPool, err = ants.NewPool(Opts.PoolSize)
	}
	return err
}

//提交到协程池处理消息
func SubmitMsgToAntsPool(c Controller, s Session, msg interface{}) {
	if err := antsPool.Submit(func() {
		c.OnMsg(s, msg)
	}); err != nil {
		Log.Errorf("antsPool commit message error,reason is ", err.Error())
	}
}
