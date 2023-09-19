package main

import (
	"fmt"
	"github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/demo/handler"
	"github.com/zjllib/gonet/v3/demo/handler/proto"
	"github.com/zjllib/gonet/v3/transport/ws"
	"log"
	"time"
)

const (
	MTU         = 1500                      // 最大传输单元
	packetLen   = 2                         // 包体大小字段
	msgIDOffset = 4                         // 消息ID字段
	headerSize  = msgIDOffset + msgIDOffset //包头部分
)

// ws://47.57.65.221:8088/game/blockInfo
// ws://192.168.0.125:8088/game/blockInfo
func main() {
	newConnection()
}

func newConnection() {
	ctx := gonet.NewContext(handler.InitClientRouter)
	session, err := ws.NewClient(ctx).Dial("ws://localhost:8088/center/ws")
	if err != nil {
		log.Fatal(err)
	}
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
