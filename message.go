package gonet

type MessageID uint32

var (
	msgNewConn = &Message{
		ID: MsgIDNewConnection,
	}
	msgConnClose = &Message{
		ID: MsgIDConnClose,
	}
)

//系统消息
const (
	MsgIDNewConnection MessageID = iota + 1
	MsgIDConnClose
)

//消息体
type Message struct {
	ID   MessageID   `json:"id"`
	Body interface{} `json:"data"`
}
