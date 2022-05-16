package main

import (
	"fmt"
	"github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/demo/proto"
	_ "github.com/zjllib/gonet/v3/transport/ws" //协议
	"log"
)

func init() {
	//消息路由
	gonet.Route(gonet.NewConnection, nil, Handler)
	gonet.Route(gonet.SessionClose, nil, Handler)
	gonet.Route(101, proto.Say{}, Handler)
}

func main() {
	server := gonet.NewServer(
		gonet.Address("ws://localhost:8088/center/ws"), //listen addr
		gonet.MaxWorkerPoolSize(20))
	log.Printf("server listening on %s", server.Addr())
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func Handler(s *gonet.Session) {
	if s.Msg != nil {
		switch s.Msg.ID {
		case gonet.NewConnection:
			log.Println("connected session_id:", s.Connection.ID(), " ip:", s.Connection.RemoteAddr().String())
		case gonet.SessionClose:
			log.Println("connected session_id:", s.Connection.ID(), " error:", s.Msg.Body)
		case 101:
			fmt.Println("session_id:", s.Connection.ID(), " say ", s.Msg.Body.(*proto.Say).Content)
			//fmt.Println(reflect.TypeOf(msg.Body))
		default:
			log.Println("unknown session_id:", s.ID)
		}
	}

}
