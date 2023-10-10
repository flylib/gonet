package main

import (
	"github.com/flylib/gonet/demo/handler"
	"github.com/flylib/gonet/transport/ws" //协议
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
