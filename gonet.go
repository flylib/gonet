package gonet

import (
	"log"
)

type MessageHandler func(IMessage)

func init() {
	log.SetPrefix("[gonet]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}
