package codec

import . "github.com/Quantumoffices/goNet"

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
func decodeMessage(msgIdx int, data []byte) (interface{}, error) {
	msg := GetMsgByIdx(msgIdx)
	err := defaultCodec.Decode(data, msg)
	return msg, err
}
