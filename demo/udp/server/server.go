package main

import (
	"github.com/zjllib/goNet"
	_ "github.com/zjllib/goNet/codec/json"
	_ "github.com/zjllib/goNet/peer/udp"
)

func main() {
	p := goNet.NewPeer(goNet.WithAddr(""))
	p.Start()
}
