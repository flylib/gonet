package gonet

type MessageID uint32

//系统消息
const (
	NewConnection MessageID = iota + 1
	SessionClose
)

//消息体
type Message struct {
	ID   MessageID   `json:"id"`
	Body interface{} `json:"data"`
}
