package main

import (
	"goNet"
	_ "goNet/codec/json"
	_ "goNet/peer/udp"
)

func main() {
	p := goNet.NewPeer("server", ":8087")
	p.Start()
}
