package binary

import (
	"github.com/davyxu/goobjfmt"
	"github.com/zjllib/gonet/v3"
)

type binaryCodec struct {
}

func (b *binaryCodec) Type() string {
	return "binary"
}

func (b *binaryCodec) Encode(msgObj interface{}) (data []byte, err error) {

	return goobjfmt.BinaryWrite(msgObj)

}

func (b *binaryCodec) Decode(data []byte, msgObj interface{}) error {

	return goobjfmt.BinaryRead(data, msgObj)
}

func init() {
	codec.SetDefaultCodec(&binaryCodec{})
}
