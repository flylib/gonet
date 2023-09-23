package json

import (
	"github.com/json-iterator/go" //高性能json编码库
	"github.com/zjllib/gonet/v3"
)

var (
	_    gonet.ICodec = new(Codec)
	json              = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Codec struct {
}

// 将结构体编码为JSON的字节数组
func (j *Codec) Encode(v interface{}) (data []byte, err error) {
	return json.Marshal(v)

}

// 将JSON的字节数组解码为结构体
func (j *Codec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
