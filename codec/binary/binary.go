package binary

import (
	"github.com/davyxu/goobjfmt"
)

type BinaryCodec struct {
}

func (b BinaryCodec) Type() string {
	return "binary"
}

func (b BinaryCodec) Encode(msgObj interface{}) (data []byte, err error) {
	return goobjfmt.BinaryWrite(msgObj)

}

func (b BinaryCodec) Decode(data []byte, msgObj interface{}) error {
	return goobjfmt.BinaryRead(data, msgObj)
}
