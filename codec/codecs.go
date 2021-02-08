package codec

import (
	"fmt"
	. "github.com/zjllib/goNet"
)

type Codec interface {
	Encode(v interface{}) (data []byte, err error)
	Decode(data []byte, vObj interface{}) error
	Type() string
}

var (
	defaultCodec Codec
)

//默认编码器
func SetDefaultCodec(codec Codec) {
	defaultCodec = codec
}

//编码消息
func encodeMessage(msg interface{}) ([]byte, error) {
	return defaultCodec.Encode(msg)
}

// 解码消息
func decodeMessage(msgID uint32, data []byte) (interface{}, error) {
	msg := GetMsg(msgID)
	if msg == nil {
		return nil, fmt.Errorf("msg_%d not found", msgID)
	}
	err := defaultCodec.Decode(data, msg)
	return msg, err
}
