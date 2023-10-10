package handler

import (
	"fmt"
	"github.com/flylib/gonet/demo/handler/proto"
	"log"
)

// 消息路由
func InitClientRouter(ctx *gonet.AppContext) error {
	ctx.Route(gonet.MessageID_SessionConnect, nil, clientHandler)
	ctx.Route(gonet.MessageID_SessionClose, nil, clientHandler)
	ctx.Route(101, proto.Say{}, clientHandler)
	return nil
}

func clientHandler(msg gonet.IMessage) {
	s := msg.FromSession()
	switch msg.ID() {
	case gonet.MessageID_SessionConnect:
		log.Println("connected session_id:", s.ID(), " ip:", s.RemoteAddr().String())
	case gonet.MessageID_SessionClose:
		log.Println("connected session_id:", s.ID(), " error:", msg.Body())
	case 101:
		fmt.Println("session_id:", s.ID(), " say ", msg.Body().(*proto.Say).Content)
	default:
		log.Println("unknown message id:", msg.ID())
	}
}
