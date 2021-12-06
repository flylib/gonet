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
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func Handler(msg *gonet.Message) {
	switch msg.ID {
	case gonet.SessionConnect:
	case gonet.SessionClose:
	default:

	}
}
