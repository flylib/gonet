package main

import (
	"fmt"
	"github.com/zjllib/goNet"
	_ "github.com/zjllib/goNet/codec/json"
	_ "github.com/zjllib/goNet/peer/udp"
	"time"
)

func main() {
	p := goNet.NewPeer(goNet.WithAddr(":88"))
	p.Start()
	fmt.Println("session count=", goNet.sessions.GetSessionCount())
	s := goNet.sessions.FindSession(1)
	fmt.Println(s.ID())
	for {
		time.Sleep(time.Second)
		s.Send(goNet.Ping{TimeStamp: time.Now().Unix()})
	}
}
