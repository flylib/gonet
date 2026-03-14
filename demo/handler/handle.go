package handler

import (
	"fmt"

	"github.com/flylib/gonet"
)

type EventHandler struct{}

func (h EventHandler) OnConnect(session gonet.ISession) {
	fmt.Printf("[connect] session=%d remote=%s\n", session.ID(), session.RemoteAddr())
}

func (h EventHandler) OnClose(session gonet.ISession, err error) {
	fmt.Printf("[close] session=%d err=%v\n", session.ID(), err)
}

func (h EventHandler) OnMessage(msg gonet.IMessage) {
	switch msg.ID() {
	case 101:
		fmt.Printf("[msg] session=%d say: %s\n", msg.From().ID(), string(msg.Body()))
	}
}

func (h EventHandler) OnError(session gonet.ISession, err error) {
	fmt.Printf("[error] session=%d err=%v\n", session.ID(), err)
}
