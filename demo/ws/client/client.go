package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/zjllib/gonet"
	"github.com/zjllib/gonet/demo/ws/proto"
	_ "github.com/zjllib/gonet/peer/ws"
	_ "github.com/zjllib/gonet/v3/json"
	"time"
)

func init() {
	gonet.RegisterMsg(0, proto.MsgIDPing, proto.Ping{})
	gonet.RegisterMsg(0, proto.MsgIDPong, proto.Pong{})
}

//ws://47.57.65.221:8088/game/blockInfo
//ws://192.168.0.125:8088/game/blockInfo
func main() {
	client := gonet.NewClient("ws://localhost:8088/center/ws")
	client.Start()
	for {
		session, ok := gonet.GetSession(uint64(gonet.SessionCount()))
		if ok {
			for {
				session.Send(proto.Ping{At: time.Now().Unix()})
				logs.Info("-------- ")
				//time.Sleep(time.Millisecond)
			}
		}
		time.Sleep(time.Second)
	}
}
