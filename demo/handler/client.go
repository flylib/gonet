package handler

import (
	"fmt"
	"github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/demo/handler/proto"
	"log"
)

// 消息路由
func InitClientRouter(ctx *gonet.Context) error {
	ctx.Route(gonet.SessionConnect, nil, clientHandler)
	ctx.Route(gonet.SessionClose, nil, clientHandler)
	ctx.Route(101, proto.Say{}, clientHandler)
	return nil
}

func clientHandler(s gonet.ISession, msg gonet.IMessage) {
	switch msg.ID() {
	case gonet.SessionConnect:
		log.Println("connected session_id:", s.ID(), " ip:", s.RemoteAddr().String())
	case gonet.SessionClose:
		log.Println("connected session_id:", s.ID(), " error:", msg.Body())
	case 101:
		fmt.Println("session_id:", s.ID(), " say ", msg.Body().(*proto.Say).Content)
	default:
		log.Println("unknown message id:", msg.ID())
	}
}
