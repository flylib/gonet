package gonet

import (
	"log"
)

//var (
//	goNetContext Context //上下文
//)

func init() {
	log.SetPrefix("[gonet]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

type Hook func(msg *Message)
