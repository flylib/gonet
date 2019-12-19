package json

import (
	"encoding/json"
	"goNet/codec"
)

type jsonCodec struct {
}

// 编码器的名称
func (j *jsonCodec) Type() string {
	return "json"
}

// 将结构体编码为JSON的字节数组
func (j *jsonCodec) Encode(v interface{}) (data []byte, err error) {

	return json.Marshal(v)

}

// 将JSON的字节数组解码为结构体
func (j *jsonCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func init() {
	codec.SetDefaultCodec(&jsonCodec{})
}
