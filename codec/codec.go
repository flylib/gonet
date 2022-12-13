package codec

const (
	Json     = "json"
	Binary   = "binary"
	Protobuf = "protobuf"
	Xml      = "xml"
)

type Codec interface {
	Encode(v interface{}) (data []byte, err error)
	Decode(data []byte, vObj interface{}) error
	Type() string
}
