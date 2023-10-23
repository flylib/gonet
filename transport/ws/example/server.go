package main

import (
	"fmt"
	"github.com/flylib/gonet"
	"github.com/flylib/gonet/demo/handler"
	"github.com/flylib/gonet/transport/ws"
	"log"
)

func main() {
	ctx := gonet.NewContext(
		gonet.WithMessageHandler(handler.MessageHandler),
	)
	fmt.Println("server listen on ws://localhost:8088/center/ws")
	if err := ws.NewServer(ctx).Listen("ws://localhost:8088/center/ws"); err != nil {
		log.Fatal(err)
	}
}
