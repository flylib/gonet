package main

import (
	"github.com/flylib/gonet"
	"github.com/flylib/gonet/transport/ws"
	"log"
)

func main() {
	ctx := gonet.NewContext(
		gonet.WithPoolMaxRoutines(20),
	)
	if err := ws.NewServer(ctx).Listen("ws://localhost:8088/center/ws"); err != nil {
		log.Fatal(err)
	}
}
