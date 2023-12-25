package handler

import (
	"fmt"
	"github.com/flylib/gonet"
)

type EventHandler struct {
}

func (e EventHandler) OnError(session gonet.ISession, err error) {
	println(fmt.Sprintf("sesson-%d error-%v", session, err))
}

func (e EventHandler) OnConnect(session gonet.ISession) {
	fmt.Println(fmt.Sprintf("new session-%d from-%s", session.ID(), session.RemoteAddr()))
}

func (e EventHandler) OnClose(session gonet.ISession, err error) {
	//TODO implement me
	fmt.Println(fmt.Sprintf("session close-%d", session.ID()))
}

func (e EventHandler) OnMessage(message gonet.IMessage) {
	//TODO implement me
	switch message.ID() {
	case 101:
		fmt.Println("session-", message.From().ID(), " say:", string(message.Body()))
	}
}
