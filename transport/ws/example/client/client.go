package main

import (
	"demo/proto"
	"fmt"
	"github.com/flylib/gonet"
	"github.com/flylib/gonet/demo/handler"
	"github.com/flylib/gonet/transport/ws"
	"log"
	"time"
)

func main() {
	ctx := gonet.NewContext(
		gonet.WithMessageHandler(handler.MessageHandler),
	)
	session, err := ws.NewClient(ctx, ws.HandshakeTimeout(5*time.Second)).Dial("ws://localhost:8088/center/ws")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connect success")

	tick := time.Tick(time.Second * 3)
	var i int
	for range tick {
		fmt.Println("send msg", i)
		i++
		err = session.Send(101, &proto.Say{
			fmt.Sprintf("hello server %d", i),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
