package main

import (
	"fmt"
	"github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/demo/proto"
	"github.com/zjllib/gonet/v3/transport/gnet" //协议
	"log"
)

var context *gonet.Context

func main() {
	context = gonet.NewContext(
		gonet.Server(gnet.NewServer("tcp://:9001")),
		gonet.MaxWorkerPoolSize(20))
	InitRouter(context)
	println("server listen on:", context.Server().Addr())
	if err := context.Server().Listen(); err != nil {
		log.Fatal(err)
	}
}

// 消息路由
func InitRouter(c *gonet.Context) {
	c.Route(gonet.SessionConnect, nil, Handler)
	c.Route(gonet.SessionClose, nil, Handler)
	c.Route(101, proto.Say{}, Handler)
}

func Handler(s gonet.ISession, msg gonet.IMessage) {
	fmt.Println("cur session count:", context.SessionCount())
	switch msg.ID() {
	case gonet.SessionConnect:
		log.Println("connected session_id:", s.ID(), " ip:", s.RemoteAddr().String())
	case gonet.SessionClose:
		log.Println("connected session_id:", s.ID(), " error:", msg.Body())
	case 101:
		fmt.Println("session_id:", s.ID(), " say ", msg.Body().(*proto.Say).Content)
		//fmt.Println(reflect.TypeOf(msg.Body))
	default:
		log.Println("unknown message id:", msg.ID())
	}
}
