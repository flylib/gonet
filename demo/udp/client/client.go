package main

import (
	"fmt"
	"github.com/zjllib/gonet"
	_ "github.com/zjllib/gonet/peer/udp"
	_ "github.com/zjllib/gonet/v3/json"
	"time"
)

func main() {
	p := gonet.NewPeer(gonet.WithAddr(":88"))
	p.Start()
	fmt.Println("session count=", gonet.sessions.GetSessionCount())
	s := gonet.sessions.FindSession(1)
	fmt.Println(s.ID())
	for {
		time.Sleep(time.Second)
		s.Send(gonet.Ping{TimeStamp: time.Now().Unix()})
	}
}
