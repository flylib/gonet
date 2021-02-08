package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/zjllib/goNet"
	_ "github.com/zjllib/goNet/codec/json"
	"github.com/zjllib/goNet/demo/ws/proto"
	_ "github.com/zjllib/goNet/peer/ws"
	"time"
)

func init() {
	goNet.RegisterMsg(0, proto.MsgIDPing, proto.Ping{})
	goNet.RegisterMsg(0, proto.MsgIDPong, proto.Pong{})
}

//ws://47.57.65.221:8088/game/blockInfo
//ws://192.168.0.125:8088/game/blockInfo
func main() {
	client := goNet.NewClient("ws://localhost:8088/center/ws")
	client.Start()
	for {
		session, ok := goNet.GetSession(uint64(goNet.SessionCount()))
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
