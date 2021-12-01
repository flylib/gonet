package gonet

var (
	msgSessionConnect = SessionConnect{}
	msgSessionClose   = SessionClose{}
)

type MessageID uint32

//系统消息
const (
	MsgIDDecPoolSize MessageID = iota
	MsgIDSessionConnect
	MsgIDSessionClose
)

//消息体
type Message struct {
	Session
	ID   MessageID   `json:"id"`
	Body interface{} `json:"data"`
}

//会话
type SessionConnect struct {
}
type SessionClose struct {
}
