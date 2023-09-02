package main

import (
	"fmt"
	"github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/demo/proto"
	"github.com/zjllib/gonet/v3/transport/ws" //协议
	"log"
)

func main() {
	context := gonet.NewContext(
		gonet.Server(ws.NewServer("ws://localhost:8088/center/ws")),
		gonet.MaxWorkerPoolSize(20))
	InitRouter(context)
	println("server listen on:", context.Server().Addr())
	if err := context.Server().Listen(); err != nil {
		log.Fatal(err)
	}
}

func InitRouter(c *gonet.Context) {
	//消息路由
	c.Route(gonet.SessionConnect, nil, Handler)
	c.Route(gonet.SessionClose, nil, Handler)
	c.Route(101, proto.Say{}, Handler)
}

func Handler(msg *gonet.Message) {
	switch msg.ID {
	case gonet.SessionConnect:
		log.Println("connected session_id:", msg.GetSession().ID(), " ip:", msg.GetSession().RemoteAddr().String())
	case gonet.SessionClose:
		log.Println("connected session_id:", msg.GetSession().ID(), " error:", msg.Body)
	case 101:
		fmt.Println("session_id:", msg.GetSession().ID(), " say ", msg.Body.(*proto.Say).Content)
		//fmt.Println(reflect.TypeOf(msg.Body))
	default:
		log.Println("unknown message id:", msg.ID)
	}
}
