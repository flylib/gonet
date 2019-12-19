package main

import (
	"fmt"
	"goNet"
	_ "goNet/codec/json"
	_ "goNet/peer/udp"
	"time"
)

func main() {
	p := goNet.NewPeer("client", "127.0.0.1:8087")
	p.Start()
	fmt.Println("session count=", goNet.SessionManager.GetSessionCount())
	s := goNet.SessionManager.GetSessionById(1)
	fmt.Println(s.ID())
	for {
		time.Sleep(time.Second)
		s.Send(goNet.Ping{TimeStamp: time.Now().Unix()})
	}
}
