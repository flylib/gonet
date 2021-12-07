package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

//ws://47.57.65.221:8088/game/blockInfo
//ws://192.168.0.125:8088/game/blockInfo
func main() {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}
	conn, _, err := dialer.Dial("ws://localhost:8088/center/ws", nil)
	if err != nil {
		fmt.Println(err)
	}
	for {
		fmt.Println("start read msg")
		messageType, bytes, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
		}
		switch messageType {
		case websocket.TextMessage:
			fmt.Println(string(bytes))
		case websocket.BinaryMessage:
		case websocket.CloseMessage:
			fmt.Println("remote server closed")
			err := conn.Close()
			if err != nil {
				fmt.Println(err)
			}
		case websocket.PingMessage:
			fmt.Println("ping at ", time.Now())
		case websocket.PongMessage:
			fmt.Println("ping at ", time.Now())
		default:
			fmt.Println("unknown msg type ", messageType)
		}
	}
}
