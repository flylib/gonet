package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zjllib/gonet/v3/demo/proto"
	"net/http"
	"time"
)

const (
	MTU         = 1500                      // 最大传输单元
	packetLen   = 2                         // 包体大小字段
	msgIDOffset = 4                         // 消息ID字段
	headerSize  = msgIDOffset + msgIDOffset //包头部分
)

//ws://47.57.65.221:8088/game/blockInfo
//ws://192.168.0.125:8088/game/blockInfo
func main() {
	for {
		time.Sleep(time.Second * 10)
		go test()
	}
}

func test() {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}
	conn, _, err := dialer.Dial("ws://localhost:8088/center/ws", nil)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		tick := time.Tick(time.Second * 10)
		for {
			<-tick

			msg := proto.Say{
				Content: "hi,I`m gonet for test msg",
			}
			fmt.Println("send msg ", msg)
			arrBytes, _ := json.Marshal(msg)
			pktData := make([]byte, msgIDOffset, msgIDOffset+len(arrBytes))
			binary.LittleEndian.PutUint32(pktData, uint32(101))
			pktData = append(pktData, arrBytes...)
			conn.WriteMessage(websocket.TextMessage, pktData)
		}
	}()

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
