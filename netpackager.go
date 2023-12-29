package gonet

import "encoding/binary"

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

// 网络包解析器(network package)
type INetPackager interface {
	Package(s ISession, msgID uint32, v any) ([]byte, error)
	UnPackage(s ISession, data []byte) (Message, int, error)
}

type DefaultNetPackager struct {
}

func (d *DefaultNetPackager) Package(s ISession, msgID uint32, v any) ([]byte, error) {
	body, err := s.GetContext().Marshal(v)
	if err != nil {
		return nil, err
	}
	content := make([]byte, MsgIDOffset+len(body))
	binary.LittleEndian.PutUint32(content, msgID)
	copy(content[MsgIDOffset:], body)
	return content, nil
}

func (d *DefaultNetPackager) UnPackage(s ISession, data []byte) (*Message, int, error) {
	msgID := binary.LittleEndian.Uint32(data[:MsgIDOffset])
	return &Message{id: msgID, body: data[MsgIDOffset:], session: s}, 0, nil
}
