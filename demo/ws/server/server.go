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
			Addr:     "ws://192.168.0.125:8083/center/ws",
			PeerType: goNet.PeertypeServer,
		})
	go p.Start()
	for {
		time.Sleep(time.Second * 6)
		session, ok := goNet.FindSession(uint64(goNet.SessionCount()))
		if ok {
			session.Close()
		}
		time.Sleep(time.Minute)
	}
}
