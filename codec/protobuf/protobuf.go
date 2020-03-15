package gogopb

import (
	"github.com/Quantumoffices/goNet/codec"
	"github.com/gogo/protobuf/proto"
)

type protobufCodec struct {
}

// 编码器的名称
func (g *protobufCodec) Type() string {
	return "protobuf"
}

func (g *protobufCodec) Encode(msgObj interface{}) (data []byte, err error) {

	return proto.Marshal(msgObj.(proto.Message))

}

func (g *protobufCodec) Decode(data []byte, msgObj interface{}) error {

	return proto.Unmarshal(data, msgObj.(proto.Message))
}

func init() {
	codec.SetDefaultCodec(&protobufCodec{})
}
