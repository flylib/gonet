package gonet

import "encoding/binary"

//++++++++++++++++++++++++++++++++++++++++++++++++++++++
//+                   +              +                 +
//+  消息总长度（2）    + 消息ID（4）   + 消息内容         +
//+                   +              +                 +
//++++++++++++++++++++++++++++++++++++++++++++++++++++++

const (
	MTU        = 1500                  // 最大传输单元
	PktSizeLen = 2                     // 包体大小字段
	MsgIDLen   = 4                     // 消息ID字段
	HeaderLen  = PktSizeLen + MsgIDLen //包头部分
)

// 网络包解析器(network package)
type INetPackager interface {
	Package(msgID uint32, v any) ([]byte, error)
	UnPackage(s ISession, data []byte) (*Message, error)
}

type DefaultNetPackager struct {
}

func (d DefaultNetPackager) Package(msgID uint32, v any) ([]byte, error) {
	content, err := defaultCtx.codec.Marshal(v)
	if err != nil {
		return nil, err
	}
	body := make([]byte, MsgIDLen+len(content))
	binary.LittleEndian.PutUint32(body, msgID)
	copy(body[MsgIDLen:], content)
	return body, nil
}

func (d DefaultNetPackager) UnPackage(s ISession, data []byte) (*Message, error) {
	msgID := binary.LittleEndian.Uint32(data[:MsgIDLen])
	return &Message{id: msgID, body: data[MsgIDLen:], session: s}, nil
}

type TcpNetPackager struct {
}

func (d TcpNetPackager) Package(msgID uint32, v any) ([]byte, error) {
	content, err := defaultCtx.codec.Marshal(v)
	if err != nil {
		return nil, err
	}
	bodySize := MsgIDLen + len(content)
	buf := make([]byte, PktSizeLen+bodySize)
	binary.LittleEndian.PutUint16(buf, uint16(bodySize))
	binary.LittleEndian.PutUint32(buf, msgID)
	copy(buf[MsgIDLen:], content)
	return buf, nil
}

func (d TcpNetPackager) UnPackage(s ISession, data []byte) (*Message, error) {
	msgID := binary.LittleEndian.Uint32(data[:MsgIDLen])
	return &Message{id: msgID, body: data[MsgIDLen:], session: s}, nil
}
