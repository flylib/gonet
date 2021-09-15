package main

import (
	"github.com/zjllib/gonet"
	_ "github.com/zjllib/gonet/peer/tcp"
	_ "github.com/zjllib/gonet/v3/json"
)

func main() {
	p := gonet.NewPeer(
		gonet.WithPeerType(gonet.PeertypeServer),
		gonet.WithAddr(":8087"),
	)
	p.Start()
}
