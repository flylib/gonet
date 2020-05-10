package main

import (
	"github.com/Quantumoffices/goNet"
	_ "github.com/Quantumoffices/goNet/codec/json"
	_ "github.com/Quantumoffices/goNet/peer/ws"
)

func main() {
	p := goNet.NewPeer(
		goNet.Options{
			Addr:          "ws://127.0.0.1:8085/echo",
			PeerType:      goNet.PEERTYPE_SERVER,
			ReadDeadline:  0,
			WriteDeadline: 0,
			PoolSize:      0,
			PanicHandler:  nil,
			AllowMaxConn:  0,
		})
	p.Start()
}
