package main

import (
	"goNet"
	_ "goNet/codec/json"
	_ "goNet/peer/ws"
)

func main() {
	p := goNet.NewPeer("server", "ws://:8087/echo")
	p.Start()
}
