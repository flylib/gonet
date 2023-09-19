package handler

import (
	"fmt"
	"github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/demo/handler/proto"
	"log"
)

// 消息路由
func InitServerRouter(ctx *gonet.Context) error {
	ctx.Route(gonet.SessionConnect, nil, serverHandler)
	ctx.Route(gonet.SessionClose, nil, serverHandler)
	ctx.Route(101, proto.Say{}, serverHandler)
	return nil
}

func serverHandler(s gonet.ISession, msg gonet.IMessage) {
	switch msg.ID() {
	case gonet.SessionConnect:
		log.Println("connected session_id:", s.ID(), " ip:", s.RemoteAddr().String())
	case gonet.SessionClose:
		log.Println("connected session_id:", s.ID(), " error:", msg.Body())
	case 101:
		fmt.Println("session_id:", s.ID(), " say ", msg.Body().(*proto.Say).Content)
		//fmt.Println(reflect.TypeOf(msg.Body))
		err := s.Send(proto.Say{Content: "hell client"})
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Println("unknown message id:", msg.ID())
	}
}
