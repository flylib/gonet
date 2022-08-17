package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/zjllib/gonet/v3/demo/proto"
	"log"
	"time"
)

const (
	MTU         = 1500                    // 最大传输单元
	packetLen   = 2                       // 包体大小字段
	msgIDOffset = 4                       // 消息ID字段
	headerSize  = packetLen + msgIDOffset //包头部分
)

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	conn, err := quic.DialAddr("localhost:8088", tlsConf, nil)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	stream1, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	{
		msg := proto.Say{
			Content: "hi,I`m gonet for test msg",
		}
		fmt.Println("send msg ", msg)
		arrBytes, _ := json.Marshal(msg)
		pktData := make([]byte, headerSize, headerSize+len(arrBytes))
		binary.LittleEndian.PutUint16(pktData, uint16(headerSize+len(arrBytes)))
		binary.LittleEndian.PutUint32(pktData[packetLen:], uint32(101))
		pktData = append(pktData, arrBytes...)
		stream1.Write(pktData)
		go func() {
			time.Sleep(time.Second * 3)
			stream1.Close()
			fmt.Println(stream1.StreamID(), "close ")
		}()
	}

	go func() {
		tick := time.Tick(time.Second * 5)
		for {

			msg := proto.Say{
				Content: "hi,I`m gonet for test msg",
			}
			fmt.Println("send msg ", msg)
			arrBytes, _ := json.Marshal(msg)
			pktData := make([]byte, headerSize, headerSize+len(arrBytes))
			binary.LittleEndian.PutUint16(pktData, uint16(headerSize+len(arrBytes)))
			binary.LittleEndian.PutUint32(pktData[packetLen:], uint32(101))
			pktData = append(pktData, arrBytes...)
			stream.Write(pktData)
			<-tick
		}
	}()

	tick := time.Tick(time.Second * 1)
	for {
		<-tick
		fmt.Println("start read msg")
		var buf []byte
		n, err := stream.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(buf[:n]))
	}

}
