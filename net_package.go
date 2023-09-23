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

// message Codec
type ICodec interface {
	Encode(v interface{}) (data []byte, err error)
	Decode(data []byte, vObj interface{}) error
}

// 网络包解析器(network package)
type INetPackageParser interface {
	Package(v any) ([]byte, error)
	UnPackage(s ISession, data []byte) (IMessage, int, error)
}

type DefaultNetPackageParser struct {
	*Context
}

func (d *DefaultNetPackageParser) Package(v any) ([]byte, error) {
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

func (d *DefaultNetPackageParser) UnPackage(from ISession, data []byte) (IMessage, int, error) {
	msgID := MessageID(binary.LittleEndian.Uint32(data[:MsgIDOffset]))
	newMsg := d.Context.CreateMsg(msgID)
	err := d.Context.DecodeMessage(newMsg, data[MsgIDOffset:])
	if err != nil {
		return nil, 0, err
	}
	return &Message{id: msgID, body: newMsg, session: from}, 0, err
}
