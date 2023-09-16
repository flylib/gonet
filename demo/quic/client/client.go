package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/zjllib/gonet/v3/demo/proto"
	"log"
	"time"
)

const (
	MTU         = 1500                      // 最大传输单元
	packetLen   = 2                         // 包体大小字段
	msgIDOffset = 4                         // 消息ID字段
	headerSize  = msgIDOffset + msgIDOffset //包头部分
)

// ws://47.57.65.221:8088/game/blockInfo
// ws://192.168.0.125:8088/game/blockInfo
func main() {
	test()
	return
	for {
		time.Sleep(time.Second * 1)
		go test()
	}
}

func test() {
	//conn, err := net.Dial("tcp", "127.0.0.1:9001")
	//if err != nil {
	//	fmt.Printf("dial failed, err: %v\n", err)
	//	return
	//}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	conn, err := quic.DialAddr(context.Background(), "127.0.0.1:9001", tlsConf, nil)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		tick := time.Tick(time.Second * 3)
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
			n, err := stream.Write(pktData)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(n)
		}
	}()

	for {
		fmt.Println("start read msg")
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(buf[:n]))
	}
}
