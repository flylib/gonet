package xml

import (
	"encoding/xml"
)

type Codec struct {
}

// 将结构体编码为xml的字节数组
func (x *Codec) Encode(v interface{}) (data []byte, err error) {
	return xml.Marshal(v)

}

// 将xml的字节数组解码为结构体
func (x *Codec) Decode(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}
