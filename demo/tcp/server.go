package main

import (
	"fmt"
	"github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/demo/handler/proto"
	"github.com/zjllib/gonet/v3/transport/tcp" //协议
	"log"
)

func main() {
	context := gonet.NewContext(
		gonet.Server(tcp.NewServer(":9001")),
		gonet.WorkerPoolMaxSize(20))
	InitRouter(context)
	println("server listen on:", context.Server().Addr())
	if err := context.Server().Listen(); err != nil {
		log.Fatal(err)
	}
}

func InitRouter(c *gonet.AppContext) {
	//消息路由
	c.Route(gonet.MessageID_SessionConnect, nil, Handler)
	c.Route(gonet.MessageID_SessionClose, nil, Handler)
	c.Route(101, proto.Say{}, Handler)
}

func Handler(s gonet.ISession, msg gonet.IMessage) {
	switch msg.ID() {
	case gonet.MessageID_SessionConnect:
		log.Println("connected session_id:", s.ID(), " ip:", s.RemoteAddr().String())
	case gonet.MessageID_SessionClose:
		log.Println("connected session_id:", s.ID(), " error:", msg.Body())
	case 101:
		fmt.Println("session_id:", s.ID(), " say ", msg.Body().(*proto.Say).Content)
		//fmt.Println(reflect.TypeOf(msg.Body))
	default:
		log.Println("unknown message id:", msg.ID())
	}
}
