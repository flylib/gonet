package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/zjllib/gonet"
	"github.com/zjllib/gonet/demo/ws/proto"
	_ "github.com/zjllib/gonet/peer/ws"
	_ "github.com/zjllib/gonet/v3/json"
	"time"
)

var loginScene server

const (
	SceneLogin uint8 = 1
)

func init() {
	//登录场景
	gonet.AddCommonScene(SceneLogin, loginScene)
	gonet.RegisterMsg(SceneLogin, gonet.MsgIDSessionConnect, gonet.SessionConnect{})
	gonet.RegisterMsg(SceneLogin, gonet.MsgIDSessionClose, gonet.SessionClose{})
	gonet.RegisterMsg(SceneLogin, proto.MsgIDPing, proto.Ping{})
	gonet.RegisterMsg(SceneLogin, proto.MsgIDPong, proto.Pong{})
}

func main() {
	server := gonet.NewServer("ws://localhost:8088/center/ws")
	server.Start()
}

type server struct {
}

func (server) Handler(msg *gonet.Msg) {
	switch data := msg.Data.(type) {
	case *gonet.SessionConnect:
		logs.Info("session_%d connected at %v", msg.Session.ID(), time.Now())
	case *gonet.SessionClose:
		logs.Warn("session_%d close at %v", msg.Session.ID(), time.Now())
	case *proto.Ping:
		logs.Info("session_%d ping at %d", msg.Session.ID(), data.At)
	}
}
