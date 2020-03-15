package main

import (
	"github.com/Quantumoffices/goNet"
	_ "github.com/Quantumoffices/goNet/codec/json"
	_ "github.com/Quantumoffices/goNet/peer/ws"
)

func main() {
	p := goNet.NewPeer(
		goNet.WithPeerType(goNet.PEERTYPE_SERVER),
		goNet.WithAddr("ws://127.0.0.1:8085/echo"),
		//goNet.WithAddr("ws://:8085/echo"),
	)
	p.Start()
}
