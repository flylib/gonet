package main

import (
	"fmt"
	"github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/demo/proto"
	_ "github.com/zjllib/gonet/v3/transport/ws"
	"log"
)

func init() {
	gonet.RegisterMsg(gonet.SessionConnect, nil, Handler)
	gonet.RegisterMsg(gonet.SessionClose, nil, Handler)
	gonet.RegisterMsg(101, proto.Say{}, Handler)
}

func main() {
	server := gonet.NewServer(
		gonet.Address("ws://localhost:8088/center/ws"),
		gonet.MaxWorkerPoolSize(20))
	log.Printf("server listening on %s", server.Addr())
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func Handler(msg *gonet.Message) {
	switch msg.ID {
	case gonet.SessionConnect:
		log.Println("connected session_id:", msg.Session.ID(), " ip:", msg.Session.RemoteAddr().String())
	case gonet.SessionClose:
		log.Println("connected session_id:", msg.Session.ID(), " error:", msg.Body)
	case 101:
		fmt.Println("session_id:", msg.Session.ID(), " say ", msg.Body.(*proto.Say).Content)
		//fmt.Println(reflect.TypeOf(msg.Body))
	default:
		log.Println("unknown session_id:", msg.ID)
	}
}
