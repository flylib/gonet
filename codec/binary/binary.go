package binary

import (
	"github.com/davyxu/goobjfmt"
	"github.com/zjllib/gonet/v3"
)

var (
	_ gonet.ICodec = new(Codec)
)

type Codec struct {
}

func (b *Codec) Encode(msgObj interface{}) (data []byte, err error) {
	return goobjfmt.BinaryWrite(msgObj)

}

func (b *Codec) Decode(data []byte, msgObj interface{}) error {
	return goobjfmt.BinaryRead(data, msgObj)
}
