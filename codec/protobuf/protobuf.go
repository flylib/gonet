package protobuf

import (
	"github.com/gogo/protobuf/proto"
)

type ProtobufCodec struct {
}

// 编码器的名称
func (g ProtobufCodec) Type() string {
	return "protobuf"
}

func (g ProtobufCodec) Encode(msgObj interface{}) (data []byte, err error) {

	return proto.Marshal(msgObj.(proto.Message))

}

func (g ProtobufCodec) Decode(data []byte, msgObj interface{}) error {

	return proto.Unmarshal(data, msgObj.(proto.Message))
}
