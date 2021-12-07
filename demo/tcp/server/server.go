package main

import (
	"github.com/zjllib/gonet/v3"
	_ "github.com/zjllib/gonet/v3/transport/tcp"
)

func main() {
	p := gonet.NewPeer(
		gonet.WithPeerType(gonet.PeertypeServer),
		gonet.WithAddr(":8087"),
	)
	p.Start()
}
