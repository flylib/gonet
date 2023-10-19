package handler

import (
	"fmt"
	"github.com/flylib/gonet"
	"github.com/flylib/gonet/demo/handler/proto"
	"log"
)

// 消息路由
func InitClientRouter(ctx *gonet.AppContext) error {
	ctx.Route(gonet.MessageID_Connection_Connect, nil, clientHandler)
	ctx.Route(gonet.MessageID_Connection_Close, nil, clientHandler)
	ctx.Route(101, proto.Say{}, clientHandler)
	return nil
}

func clientHandler(msg gonet.IMessage) {
	s := msg.FromSession()
	switch msg.ID() {
	case gonet.MessageID_Connection_Connect:
		log.Println("connected session_id:", s.ID(), " ip:", s.RemoteAddr().String())
	case gonet.MessageID_Connection_Close:
		log.Println("connected session_id:", s.ID(), " error:", msg.Body())
	case 101:
		fmt.Println("session_id:", s.ID(), " say ", msg.Body().(*proto.Say).Content)
	default:
		log.Println("unknown message id:", msg.ID())
	}
}
