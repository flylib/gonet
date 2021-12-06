package codec

import (
	"github.com/json-iterator/go" //高性能json编码库
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type JsonCodec struct {
}

// 编码器的名称
func (j JsonCodec) Type() string {
	return Json
}

// 将结构体编码为JSON的字节数组
func (j JsonCodec) Encode(v interface{}) (data []byte, err error) {
	return json.Marshal(v)

}

// 将JSON的字节数组解码为结构体
func (j JsonCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
