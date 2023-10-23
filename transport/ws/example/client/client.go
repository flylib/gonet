package main

import (
	"fmt"
	"github.com/flylib/gonet"
	"github.com/flylib/gonet/demo/handler/proto"
	"github.com/flylib/gonet/transport/ws"
	"log"
	"time"
)

func main() {
	ctx := gonet.NewContext()
	session, err := ws.NewClient(ctx, ws.HandshakeTimeout(5*time.Second)).Dial("ws://localhost:8088/center/ws")
	if err != nil {
		log.Fatal(err)
	}

	//for {
	//	fmt.Println("please input what you want to send:")
	//	var input string
	//	fmt.Scanln(&input)
	//	err = session.Send(proto.Say{
	//		input,
	//	})
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}

	tick := time.Tick(time.Second * 3)
	var i int
	for range tick {
		fmt.Println("send msg", i)
		i++
		err = session.Send(proto.Say{
			fmt.Sprintf("hello server %d", i),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
