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
	Encode(v interface{}) (data []byte, err error)
	Decode(data []byte, vObj interface{}) error
}

// 网络包解析器(network package)
type INetPackageParser interface {
	Package(ctx *AppContext, v any) ([]byte, error)
	UnPackage(ctx *AppContext, s ISession, data []byte) (IMessage, int, error)
}

type DefaultNetPackageParser struct {
}

func (d *DefaultNetPackageParser) Package(ctx *AppContext, v any) ([]byte, error) {
	body, err := ctx.EncodeMessage(v)
	if err != nil {
		return nil, err
	}
	p := make([]byte, MsgIDOffset, MsgIDOffset+len(body))
	msgID, ok := ctx.GetMsgID(v)
	if ok {
		binary.LittleEndian.PutUint32(p, uint32(msgID))
		p = append(p, body...)
		return p, nil
	}
	return nil, ErrorNotExistMsg
}

func (d *DefaultNetPackageParser) UnPackage(ctx *AppContext, from ISession, data []byte) (IMessage, int, error) {
	msgID := MessageID(binary.LittleEndian.Uint32(data[:MsgIDOffset]))
	newMsg := ctx.CreateMsg(msgID)
	err := ctx.DecodeMessage(newMsg, data[MsgIDOffset:])
	if err != nil {
		return nil, 0, err
	}
	return &Message{id: msgID, body: newMsg, session: from}, 0, err
}
