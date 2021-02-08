package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/zjllib/goNet"
	_ "github.com/zjllib/goNet/codec/json"
	"github.com/zjllib/goNet/demo/ws/proto"
	_ "github.com/zjllib/goNet/peer/ws"
	"time"
)

var loginScene server

const (
	SceneLogin uint8 = 1
)

func init() {
	//登录场景
	goNet.AddCommonScene(SceneLogin, loginScene)
	goNet.RegisterMsg(SceneLogin, goNet.MsgIDSessionConnect, goNet.SessionConnect{})
	goNet.RegisterMsg(SceneLogin, goNet.MsgIDSessionClose, goNet.SessionClose{})
	goNet.RegisterMsg(SceneLogin, proto.MsgIDPing, proto.Ping{})
	goNet.RegisterMsg(SceneLogin, proto.MsgIDPong, proto.Pong{})
}

func main() {
	server := goNet.NewServer("ws://localhost:8088/center/ws")
	server.Start()
}

type server struct {
}

func (server) Handler(msg *goNet.Msg) {
	switch data := msg.Data.(type) {
	case *goNet.SessionConnect:
		logs.Info("session_%d connected at %v", msg.Session.ID(), time.Now())
	case *goNet.SessionClose:
		logs.Warn("session_%d close at %v", msg.Session.ID(), time.Now())
	case *proto.Ping:
		logs.Info("session_%d ping at %d", msg.Session.ID(), data.At)
	}
}
