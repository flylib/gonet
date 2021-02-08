package main

import (
	"fmt"
	"github.com/zjllib/goNet"
	_ "github.com/zjllib/goNet/codec/json"
	_ "github.com/zjllib/goNet/peer/tcp"
	"time"
)

func main() {
	p := goNet.NewPeer(
		goNet.Options{
			PeerType: goNet.PEERTYPE_CLIENT,
			Addr:     ":8087",
		})
	p.Start()
	fmt.Println("session count=", goNet.sessions.GetSessionCount())
	session, exist := goNet.sessions.FindSession(1)
	if exist {
		for {
			time.Sleep(time.Second)
			session.Send(goNet.Ping{TimeStamp: time.Now().Unix()})
		}
	}
}
