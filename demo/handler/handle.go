package handler

import (
	"demo/proto"
	"fmt"
	"github.com/flylib/gonet"
	"log"
)

// 消息路由
func MessageHandler(msg gonet.IMessage) {
	s := msg.From()
	switch msg.ID() {
	case gonet.MessageID_Connection_Connect:
		log.Println("connected session_id:", s.ID(), " ip:", s.RemoteAddr().String())
	case gonet.MessageID_Connection_Close:
		log.Println("connected session_id:", s.ID(), " error:", msg.Body())
	case 101:
		pb := proto.Say{}
		err := msg.UnmarshalTo(&pb)
		if err != nil {
			panic(err)
		}
		fmt.Println("session_id:", s.ID(), " say ", pb.Content)
		err = s.Send(102, proto.Say{Content: "hell client"})
		if err != nil {
			log.Fatal(err)
		}
	case 102:
		pb := proto.Say{}
		err := msg.UnmarshalTo(&pb)
		if err != nil {
			panic(err)
		}
		fmt.Println("session_id:", s.ID(), " say ", pb.Content)
	default:
		log.Println("unknown message id:", msg.ID())
	}
}
