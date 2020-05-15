package main

import (
	"fmt"
	"github.com/Quantumoffices/goNet"
	_ "github.com/Quantumoffices/goNet/codec/json"
	_ "github.com/Quantumoffices/goNet/peer/tcp"
	"time"
)

func main() {
	p := goNet.NewPeer(
		goNet.WithPeerType(goNet.PEERTYPE_CLIENT),
		goNet.WithAddr(":8087"),
	)
	p.Start()
	fmt.Println("session count=", goNet.SessionManager.GetSessionCount())
	s := goNet.SessionManager.GetSessionById(1)
	fmt.Println(s.ID())
	for {
		time.Sleep(time.Second)
		s.Send(goNet.Ping{TimeStamp: time.Now().Unix()})
	}
}
