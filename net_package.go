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

// message Codec
type ICodec interface {
	Marshal(v any) (data []byte, err error)
	Unmarshal(data []byte, v any) error
}

// 网络包解析器(network package)
type INetPackageParser interface {
	Package(msgID MessageID, v any) ([]byte, error)
	UnPackage(data []byte) (IMessage, int, error)
}

type DefaultNetPackageParser struct {
	*AppContext
}

func (d *DefaultNetPackageParser) Package(msgID MessageID, v any) ([]byte, error) {
	body, err := d.EncodeMessage(v)
	if err != nil {
		return nil, err
	}
	payload := make([]byte, MsgIDOffset, MsgIDOffset+len(body))
	binary.LittleEndian.PutUint32(payload, uint32(msgID))
	payload = append(payload, body...)
	return nil, ErrorNotExistMsg
}

func (d *DefaultNetPackageParser) UnPackage(data []byte) (IMessage, int, error) {
	msgID := MessageID(binary.LittleEndian.Uint32(data[:MsgIDOffset]))
	return &Message{id: msgID, body: data[MsgIDOffset:]}, 0, nil
}
