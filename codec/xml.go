package codec

import (
	"encoding/xml"
)

type XmlCodec struct {
}

// 编码器的名称
func (x XmlCodec) Type() string {
	return Xml
}

// 将结构体编码为xml的字节数组
func (x XmlCodec) Encode(v interface{}) (data []byte, err error) {
	return xml.Marshal(v)

}

// 将xml的字节数组解码为结构体
func (x XmlCodec) Decode(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}
