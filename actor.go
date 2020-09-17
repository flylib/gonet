package goNet

import (
	"github.com/astaxie/beego/logs"
	"time"
)

//默认ID
const DefaultActorID = 0

//系统消息路由
var defaultActor Actor = &DefaultActor{}

//使用Actor模型的好处：
//事件模型驱动--Actor之间的通信是异步的，即使Actor在发送消息后也无需阻塞或者等待就能够处理其他事情
//强隔离性--Actor中的方法不能由外部直接调用，所有的一切都通过消息传递进行的，从而避免了Actor之间的数据共享，想要
//观察到另一个Actor的状态变化只能通过消息传递进行询问
//位置透明--无论Actor地址是在本地还是在远程机上对于代码来说都是一样的
//轻量性--Actor是非常轻量的计算单机，只需少量内存就能达到高并发
type Actor interface {
	Receive(Context)
}

//系统控制模块
type DefaultActor struct {
}

func (*DefaultActor) Receive(c Context) {
	switch data := c.Message().(type) {
	case *SessionConnect:
		logs.Info("session_%v connected", c.Session().ID())
	case *SessionClose:
		logs.Warn("session_%v closed", c.Session().ID())
	case *Ping:
		logs.Info("session_%v ping at time=%v", c.Session().ID(), time.Unix(data.TimeStamp, 0).String())
		c.Session().Send(Pong{TimeStamp: time.Now().Unix()})
	case *Pong:
		logs.Info("session_%v pong at time=%v", c.Session().ID(), time.Unix(data.TimeStamp, 0).String())
	}
}

func UpdateSysActor(c Actor) {
	if c == nil {
		return
	}
	defaultActor = c
}
