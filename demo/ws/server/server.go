package main

import (
	"goNet"
	_ "goNet/codec/json"
	_ "goNet/peer/ws"
)

func main() {
	p := goNet.NewPeer(
		goNet.WithPeerType(goNet.PEERTYPE_SERVER),
		goNet.WithAddr("ws://127.0.0.1:8085/echo"),
		//goNet.WithAddr("ws://:8085/echo"),
	)
	p.Start()
}
