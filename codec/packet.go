package codec

import (
	"encoding/binary"
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"io"
	"net"
	"reflect"
)

//++++++++++++++++++++++++++++++++++++++++++++++++++++++
//+                   +              +                 +
//+  消息总长度（2）    + 消息ID（4）   + 消息内容         +
//+                   +              +                 +
//++++++++++++++++++++++++++++++++++++++++++++++++++++++

const (
	MTU         = 1500                      // 最大传输单元
	packetLen   = 2                         // 包体大小字段
	msgIDOffset = 4                         // 消息ID字段
	headerSize  = msgIDOffset + msgIDOffset //包头部分
)

//----------------------------------------------【解析包】--------------------------------------------------
//func ParserPacket(data []byte) (int, interface{}, error) {
//	// 小于包头
//	if len(data) < packetLen {
//		return 0, nil, errors.New("packet size too min")
//	}
//	// 读取Size
//	size := binary.LittleEndian.Uint16(data)
//	// 出错，等待下次数据
//	if size > MTU {
//		return DefaultSceneID, nil, errors.New(fmt.Sprintf("packet size %v max MTU length", size))
//	}
//	// 读取消息ID
//	msgId := int(binary.LittleEndian.Uint16(data[packetLen:]))
//	//内容
//	content := data[headerSize : headerSize+size]
//	ActorID := GetMsgSceneID(msgId)
//	msg, err := decodeMessage(msgId, content)
//	return ActorID, msg, err
//}

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
	msgID, ok := GetMsgID(reflect.TypeOf(msg))
	if ok {
		// ID
		binary.LittleEndian.PutUint16(pktData[2:], uint16(msgID))
		// Value
		copy(pktData[headerSize:], body)
		_, err = w.Write(pktData)
		return err
	}
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
	msgID, ok := GetMsgID(reflect.TypeOf(msg))
	if ok {
		// ID
		binary.LittleEndian.PutUint16(pktData[2:], uint16(msgID))
		// Value
		copy(pktData[headerSize:], body)

		_, err = w.WriteToUDP(pktData, toAddr)
		return err
	}
	return err
}

//----------------------------------------------【ws】--------------------------------------------------
func ParserWSPacket(pkt []byte) (*Message, error) {
	msgID := Bytes2Uint32(pkt[:msgIDOffset])
	msg, err := decodeMessage(msgID, pkt[msgIDOffset:])
	if err != nil {
		return nil, err
	}
	return &Message{ID: MessageID(msgID), Body: msg}, err
}

func SendWSPacket(w *websocket.Conn, msg interface{}) error {
	arrBytes, err := encodeMessage(msg)
	if err != nil {
		return err
	}
	pktData := make([]byte, msgIDOffset, msgIDOffset+len(arrBytes))
	msgID, ok := GetMsgID(msg)
	if ok {
		binary.LittleEndian.PutUint32(pktData, uint32(msgID))
		pktData = append(pktData, arrBytes...)
		return w.WriteMessage(websocket.TextMessage, pktData)
	}
	return ErrorNotExistMsg
}

func Bytes2Uint32(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}
