package main

import (
	"goNet"
	_ "goNet/codec/json"
	_ "goNet/peer/tcp"
)

func main() {
	p := goNet.NewPeer(
		goNet.WithPeerType(goNet.PEER_SERVER),
		goNet.WithAddr(":8087"),
	)
	p.Start()
}
