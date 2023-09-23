package protobuf

import (
	"github.com/gogo/protobuf/proto"
	"github.com/zjllib/gonet/v3"
)

var (
	_ gonet.ICodec = new(Codec)
)

type Codec struct {
}

func (g *Codec) Encode(msgObj interface{}) (data []byte, err error) {

	return proto.Marshal(msgObj.(proto.Message))

}

func (g *Codec) Decode(data []byte, msgObj interface{}) error {

	return proto.Unmarshal(data, msgObj.(proto.Message))
}
