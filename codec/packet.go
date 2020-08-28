package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Quantumoffices/goNet"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"reflect"
	"strconv"
)

//++++++++++++++++++++++++++++++++++++++++++++++++++++++
//+                   +              +                 +
//+  消息总长度（2）  + 消息ID（2）  + 消息内容        +
//+                   +              +                 +
//++++++++++++++++++++++++++++++++++++++++++++++++++++++

const (
	MTU        = 1500                // 最大传输单元
	packetLen  = 2                   // 包体大小字段
	msgIDLen   = 2                   // 消息ID字段
	headerSize = msgIDLen + msgIDLen //包头部分
)

//----------------------------------------------【解析包】--------------------------------------------------
func ParserPacket(data []byte) (int, interface{}, error) {
	// 小于包头
	if len(data) < packetLen {
		return goNet.DefaultRouteID, nil, errors.New("packet size too min")
	}
	// 读取Size
	size := binary.LittleEndian.Uint16(data)
	// 出错，等待下次数据
	if size > MTU {
		return goNet.DefaultRouteID, nil, errors.New(fmt.Sprintf("packet size %v max MTU length", size))
	}
	// 读取消息ID
	msgId := int(binary.LittleEndian.Uint16(data[packetLen:]))
	//内容
	content := data[headerSize : headerSize+size]
	routeID := goNet.FindMsgOnRoute(msgId)
	msg, err := decodeMessage(msgId, content)
	return routeID, msg, err
}

//----------------------------------------------【发送包】--------------------------------------------------
func SendPacket(w io.Writer, msg interface{}) error {
	// 将用户数换为字节数组和消息ID
	body, err := encodeMessage(msg)
	if err != nil {
		return err
	}
	pktData := make([]byte, headerSize+len(body))
	// Size==len(body)
	binary.LittleEndian.PutUint16(pktData, uint16(len(body)))
	// ID
	binary.LittleEndian.PutUint16(pktData[2:], uint16(goNet.GetMsgID(reflect.TypeOf(msg))))
	// Value
	copy(pktData[headerSize:], body)

	_, err = w.Write(pktData)
	return err
}

func SendUdpPacket(w *net.UDPConn, msg interface{}, toAddr *net.UDPAddr) error {
	// 将用户数换为字节数组和消息ID
	body, err := encodeMessage(msg)
	if err != nil {
		return err
	}

	pktData := make([]byte, headerSize+len(body))
	// Size==len(body)
	binary.LittleEndian.PutUint16(pktData, uint16(len(body)))
	// ID
	binary.LittleEndian.PutUint16(pktData[2:], uint16(goNet.GetMsgID(reflect.TypeOf(msg))))
	// Value
	copy(pktData[headerSize:], body)

	_, err = w.WriteToUDP(pktData, toAddr)

	return err
}

//----------------------------------------------【ws】--------------------------------------------------
func ParserWSPacket(pkt []byte) (int, interface{}, error) {
	for index, d := range pkt {
		if d == '\n' {
			msgID, err := strconv.Atoi(string(pkt[:index]))
			if err != nil {
				return goNet.DefaultRouteID, nil, err
			}
			routeID := goNet.FindMsgOnRoute(msgID)
			msg, err := decodeMessage(msgID, pkt[index+1:])
			return routeID, msg, err
		}
	}
	return goNet.DefaultRouteID, nil, errors.New("parser message error.EOF")
}

func SendWSPacket(w *websocket.Conn, msg interface{}) error {
	body, err := encodeMessage(msg)
	if err != nil {
		return err
	}
	return w.WriteMessage(websocket.TextMessage,
		bytes.Join([][]byte{[]byte(strconv.Itoa(goNet.GetMsgID(reflect.TypeOf(msg)))), body}, []byte{10}))
}
