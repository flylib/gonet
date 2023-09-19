package main

import (
	"fmt"
	"github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/demo/handler/proto"
	"github.com/zjllib/gonet/v3/transport/udp" //协议
	"log"
	"time"
)

func main() {
	context := gonet.NewContext(
		gonet.Server(udp.NewServer(":9001")),
		gonet.WorkerPoolMaxSize(20))
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

func Handler(s gonet.ISession, msg gonet.IMessage) {
	switch msg.ID() {
	case gonet.SessionConnect:
		log.Println("connected session_id:", s.ID(), " ip:", s.RemoteAddr().String())
		write(s)
	case gonet.SessionClose:
		log.Println("connected session_id:", s.ID(), " error:", msg.Body())
	case 101:
		fmt.Println("session_id:", s.ID(), " say ", msg.Body().(*proto.Say).Content)
		//fmt.Println(reflect.TypeOf(msg.Body))
	default:
		log.Println("unknown message id:", msg.ID())
	}
}

func write(s gonet.ISession) {
	go func() {
		for {
			time.Sleep(time.Second)
			err := s.Send(proto.Say{"你好"})
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
}
