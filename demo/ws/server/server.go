package main

import (
	"github.com/zjllib/gonet/v3"
	_ "github.com/zjllib/gonet/v3/transport/ws"
	"log"
)

func init() {
	gonet.RegisterMsg(gonet.SessionConnect, nil, Handler)
	gonet.RegisterMsg(gonet.SessionClose, nil, Handler)
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
		log.Println("connected session_id:", msg.Session.ID())
	default:
		log.Println("unknown session_id:", msg.ID)
	}
}
