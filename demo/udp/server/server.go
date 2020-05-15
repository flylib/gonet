package main

import (
	"github.com/Quantumoffices/goNet"
	_ "github.com/Quantumoffices/goNet/codec/json"
	_ "github.com/Quantumoffices/goNet/peer/udp"
)

func main() {
	p := goNet.NewPeer(goNet.WithAddr(""))
	p.Start()
}
