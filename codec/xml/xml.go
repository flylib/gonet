package xml

import (
	"encoding/xml"
	"github.com/Quantumoffices/goNet/codec"
)

type xmlCodec struct {
}

// 编码器的名称
func (x *xmlCodec) Type() string {
	return "xml"
}

// 将结构体编码为xml的字节数组
func (j *xmlCodec) Encode(v interface{}) (data []byte, err error) {

	return xml.Marshal(v)

}

// 将xml的字节数组解码为结构体
func (j *xmlCodec) Decode(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func init() {
	codec.SetDefaultCodec(&xmlCodec{})
}
