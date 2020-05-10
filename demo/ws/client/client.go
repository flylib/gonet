package main

import (
	"github.com/Quantumoffices/goNet"
	_ "github.com/Quantumoffices/goNet/codec/json"
	_ "github.com/Quantumoffices/goNet/peer/ws"
	"time"
)

func main() {
	p := goNet.NewPeer(
		goNet.Options{
			Addr:     "ws://127.0.0.1:8085/echo",
			PeerType: goNet.PEERTYPE_CLIENT,
			//ReadDeadline:  0,
			//WriteDeadline: 0,
			//PoolSize:      0,
			//PanicHandler:  nil,
			//AllowMaxConn:  0,
		})
	p.Start()
	s := goNet.SessionManager.GetSessionById(1)
	for {
		time.Sleep(time.Second)
		s.Send(goNet.Ping{TimeStamp: time.Now().Unix()})
	}
}
