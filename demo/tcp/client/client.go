package main

import (
	"fmt"
	"github.com/zjllib/gonet"
	_ "github.com/zjllib/gonet/peer/tcp"
	_ "github.com/zjllib/gonet/v3/json"
	"time"
)

func main() {
	p := gonet.NewPeer(
		gonet.Options{
			PeerType: gonet.PEERTYPE_CLIENT,
			Addr:     ":8087",
		})
	p.Start()
	fmt.Println("session count=", gonet.sessions.GetSessionCount())
	session, exist := gonet.sessions.FindSession(1)
	if exist {
		for {
			time.Sleep(time.Second)
			session.Send(gonet.Ping{TimeStamp: time.Now().Unix()})
		}
	}
}
