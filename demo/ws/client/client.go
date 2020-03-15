package main

import (
	"fmt"
	"github.com/Quantumoffices/goNet"
	_ "github.com/Quantumoffices/goNet/codec/json"
	_ "github.com/Quantumoffices/goNet/peer/ws"
	"time"
)

func main() {
	p := goNet.NewPeer(
		goNet.WithPeerType(goNet.PEERTYPE_CLIENT),
		goNet.WithAddr("ws://127.0.0.1:8085/echo"),
		//goNet.WithAddr("ws://www.quantumstudio.cn:8000/echo"),
		//goNet.WithAddr("ws://www.quantumstudio.cn:8000/center_server_cluster"),
	)
	p.Start()
	fmt.Println("session count=", goNet.SessionManager.GetSessionCount())
	s := goNet.SessionManager.GetSessionById(1)
	//s.Close()
	for {
		time.Sleep(time.Second)
		s.Send(goNet.Ping{TimeStamp: time.Now().Unix()})
	}
}
