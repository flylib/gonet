package main

import (
	"github.com/zjllib/gonet"
	_ "github.com/zjllib/gonet/peer/udp"
	_ "github.com/zjllib/gonet/v3/json"
)

func main() {
	p := gonet.NewPeer(gonet.WithAddr(""))
	p.Start()
}
