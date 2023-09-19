package main

import (
	"github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/demo/handler"
	"github.com/zjllib/gonet/v3/transport/ws" //协议
	"log"
)

func main() {
	ctx := gonet.NewContext(
		gonet.WorkerPoolMaxSize(20),
		handler.InitServerRouter,
	)

	if err := ws.NewServer(ctx).Listen("ws://localhost:8088/center/ws"); err != nil {
		log.Fatal(err)
	}
}
