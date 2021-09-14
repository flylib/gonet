package goNet

type Conn interface {
	//原始套接字
	Socket() interface{}
	// 发送消息，消息需要以指针格式传入
	Send(msg interface{})
	// 断开
	Close()
	// ID
	ID() uint64
	//数据存储
	Value(obj ...interface{}) interface{}
	//添加场景,如果场景相同会进行覆盖
	JoinScene(sceneID uint8, scene Scene)
	//获取场景
	GetScene(sceneID uint8) Scene
}
