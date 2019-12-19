package main

import (
	"fmt"
	"goNet"
	_ "goNet/codec/json"
	_ "goNet/peer/ws"
	"time"
)

func main() {
	p := goNet.NewPeer("client", "ws://127.0.0.1:8087/echo")
	p.Start()
	fmt.Println("session count=", goNet.SessionManager.GetSessionCount())
	s := goNet.SessionManager.GetSessionById(1)
	//s.Close()
	for {
		time.Sleep(time.Second)
		s.Send(goNet.Ping{TimeStamp: time.Now().Unix()})
	}
}
