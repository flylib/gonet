package main

import (
	"github.com/zjllib/goNet"
	_ "github.com/zjllib/goNet/codec/json"
	_ "github.com/zjllib/goNet/peer/tcp"
)

func main() {
	p := goNet.NewPeer(
		goNet.WithPeerType(goNet.PeertypeServer),
		goNet.WithAddr(":8087"),
	)
	p.Start()
}
