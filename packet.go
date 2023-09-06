package gonet

import "encoding/binary"

//
//import (
//	"encoding/binary"
//	"github.com/gorilla/websocket"
//	"io"
//	"net"
//	"reflect"
//)
//
//++++++++++++++++++++++++++++++++++++++++++++++++++++++
//+                   +              +                 +
//+  消息总长度（2）    + 消息ID（4）   + 消息内容         +
//+                   +              +                 +
//++++++++++++++++++++++++++++++++++++++++++++++++++++++

const (
	MTU           = 1500                      // 最大传输单元
	PktSizeOffset = 2                         // 包体大小字段
	MsgIDOffset   = 4                         // 消息ID字段
	HeaderOffset  = MsgIDOffset + MsgIDOffset //包头部分
)

type IPackageParser interface {
	Package(v any) ([]byte, error)
	UnPackage(data []byte) (IMessage, int, error)
}

type defaultPackageParser struct {
	*Context
}

func (d *defaultPackageParser) Package(v any) ([]byte, error) {
	body, err := d.Context.EncodeMessage(v)
	if err != nil {
		return nil, err
	}
	p := make([]byte, MsgIDOffset, MsgIDOffset+len(body))
	msgID, ok := d.Context.GetMsgID(v)
	if ok {
		binary.LittleEndian.PutUint32(p, uint32(msgID))
		p = append(p, body...)
		return p, nil
	}
	return nil, ErrorNotExistMsg
}

func (d *defaultPackageParser) UnPackage(data []byte) (IMessage, int, error) {
	msgID := MessageID(Bytes2Uint32(data[:MsgIDOffset]))
	newMsg := d.Context.CreateMsg(msgID)
	err := d.Context.DecodeMessage(newMsg, data[MsgIDOffset:])
	if err != nil {
		return nil, 0, err
	}
	return &Message{id: msgID, body: newMsg}, 0, err
}

//
////----------------------------------------------【解析包】--------------------------------------------------
////用于处理tcp,udp等粘包问题
//func ParserPacket(data []byte) (*Message, int, error) {
//	dataSize := len(data)
//	// 小于包头
//	if dataSize < PktSizeOffset {
//		return nil, 0, ErrorNotExistMsg
//	}
//	//包大小
//	pktSize := binary.LittleEndian.Uint16(data)
//	// 读取消息ID
//	msgID := MessageID(binary.LittleEndian.Uint16(data[PktSizeOffset:]))
//	//内容
//	content := data[HeaderOffset:pktSize]
//
//	newMsg := CreateMsg(msgID)
//	if newMsg == nil {
//		return nil, 0, ErrorNotExistMsg
//	}
//	err := DecodeMessage(newMsg, content)
//	if err != nil {
//		return nil, 0, err
//	}
//	return &Message{ID: msgID, Body: newMsg}, dataSize - int(pktSize), err
//}
//
////----------------------------------------------【发送包】--------------------------------------------------
//func SendPacket(w io.Writer, msg interface{}) error {
//	body, err := EncodeMessage(msg)
//	if err != nil {
//		return err
//	}
//	pktSize := HeaderOffset + len(body)
//	//header+body
//	pktData := make([]byte, HeaderOffset)
//	//信息：包体大小
//	binary.LittleEndian.PutUint16(pktData, uint16(pktSize))
//	msgID, ok := GetMsgID(reflect.TypeOf(msg))
//	if ok {
//		//信息：message id
//		binary.LittleEndian.PutUint16(pktData[PktSizeOffset:], uint16(msgID))
//		// 信息：message content
//		pktData = append(pktData, body...)
//		_, err = w.Write(pktData)
//		return err
//	}
//	return ErrorNotExistMsg
//}
//
//func SendUdpPacket(w *net.UDPConn, msg interface{}, toAddr *net.UDPAddr) error {
//	// 将用户数换为字节数组和消息ID
//	body, err := EncodeMessage(msg)
//	if err != nil {
//		return err
//	}
//	pktData := make([]byte, HeaderOffset+len(body))
//	// Size==len(body)
//	binary.LittleEndian.PutUint16(pktData, uint16(len(body)))
//	msgID, ok := GetMsgID(reflect.TypeOf(msg))
//	if ok {
//		// ID
//		binary.LittleEndian.PutUint16(pktData[2:], uint16(msgID))
//		// Value
//		copy(pktData[HeaderOffset:], body)
//		_, err = w.WriteToUDP(pktData, toAddr)
//		return err
//	}
//	return err
//}
//
////----------------------------------------------【ws】--------------------------------------------------
//func ParserWSPacket(pkt []byte) (*Message, error) {
//	msgID := MessageID(Bytes2Uint32(pkt[:MsgIDOffset]))
//	newMsg := CreateMsg(msgID)
//	err := DecodeMessage(newMsg, pkt[MsgIDOffset:])
//	if err != nil {
//		return nil, err
//	}
//	return &Message{ID: msgID, Body: newMsg}, err
//}
//
//func SendWSPacket(w *websocket.Conn, msg interface{}) error {
//	arrBytes, err := EncodeMessage(msg)
//	if err != nil {
//		return err
//	}
//	pktData := make([]byte, MsgIDOffset, MsgIDOffset+len(arrBytes))
//	msgID, ok := GetMsgID(msg)
//	if ok {
//		binary.LittleEndian.PutUint32(pktData, uint32(msgID))
//		pktData = append(pktData, arrBytes...)
//		return w.WriteMessage(websocket.BinaryMessage, pktData)
//	}
//	return ErrorNotExistMsg
//}

func Bytes2Uint32(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}
