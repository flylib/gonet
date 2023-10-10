package main

import (
	"fmt"
	"github.com/flylib/gonet/demo/handler/proto"
	"github.com/flylib/gonet/transport/gnet" //协议
	"log"
)

var context *gonet.AppContext

func main() {
	context = gonet.NewContext(
		gonet.Server(gnet.NewServer("tcp://:9001")),
		gonet.WorkerPoolMaxSize(20))
	InitRouter(context)
	println("server listen on:", context.Server().Addr())
	if err := context.Server().Listen(); err != nil {
		log.Fatal(err)
	}
}

// 消息路由
func InitRouter(c *gonet.AppContext) {
	c.Route(gonet.MessageID_SessionConnect, nil, Handler)
	c.Route(gonet.MessageID_SessionClose, nil, Handler)
	c.Route(101, proto.Say{}, Handler)
}

func Handler(s gonet.ISession, msg gonet.IMessage) {
	fmt.Println("cur session count:", context.SessionCount())
	switch msg.ID() {
	case gonet.MessageID_SessionConnect:
		log.Println("connected session_id:", s.ID(), " ip:", s.RemoteAddr().String())
	case gonet.MessageID_SessionClose:
		log.Println("connected session_id:", s.ID(), " error:", msg.Body())
	case 101:
		fmt.Println("session_id:", s.ID(), " say ", msg.Body().(*proto.Say).Content)
		//fmt.Println(reflect.TypeOf(msg.Body))
		s.Send(&proto.Say{Content: "hell client"})
	default:
		log.Println("unknown message id:", msg.ID())
	}
}
