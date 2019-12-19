package main

import (
	"goNet"
	_ "goNet/codec/json"
	_ "goNet/peer/tcp"
)

func main() {
	p := goNet.NewPeer("server",":8087")
	p.Start()
}
