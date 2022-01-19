
![gonetlogo](docs/logo.jpg)
## version
 v 2.0.0
## 介绍
一个基于go语言开发的网络脚手架,参考[cellnet](https://github.com/davyxu/cellnet)和[beego](https://github.com/astaxie/beego)两大开框架的设计，使用非常方便简洁，轻松让你开发出高并发高性能的网络应用，可以用于游戏,app等任何领域的通讯。

## 主要特性及追求目标
- 高并发
- 高性能
- 简单易用
- 线性安全
- 兼容性强
- 多领域应用
- 防崩溃
- 错误快速定位

## 通讯协议支持
- [x] TCP
- [x] UDP
- [x] WEBSOCKET
- [ ] QUIC
- [ ] KCP
- [ ] HTTP
- [ ] RPC
## 数据编码格式支持
- [x] json
- [x] xml
- [x] binary
- [x] protobuf

## 关键技术
- [x] 会话池(session pool）
- [x] 协程池(goroutine pool)

## 安装教程
### **1.** git clone到 GOPATH/src目录下

```
git clone https://github.com/zjllib/gonet.git
```

## 使用样例参考
- 服务端
```go
	func main() {
	server := gonet.NewServer("ws://localhost:8088/center/ws")
	server.Start()
}
```
- 客户端
```go
	client := gonet.NewClient("ws://localhost:8088/center/ws")
	client.Start()
```
- 消息处理实现及注册
```go
//系统消息
const (
	MsgIDDecPoolSize uint32 = iota
	MsgIDSessionConnect
	MsgIDSessionClose
)
//消息注入
func init() {
	gonet.RegisterMsg(SceneLogin, gonet.MsgIDSessionConnect, gonet.SessionConnect{})
	gonet.RegisterMsg(SceneLogin, gonet.MsgIDSessionClose, gonet.SessionClose{})
	gonet.RegisterMsg(SceneLogin, proto.MsgIDPing, proto.Ping{})
	gonet.RegisterMsg(SceneLogin, proto.MsgIDPong, proto.Pong{})
}

//消息处理场景
type server struct {
}

func (server) Handler(msg *gonet.Msg) {
	switch data := msg.Data.(type) {
	case *gonet.SessionConnect:
		logs.Info("session_%d connected at %v", msg.Session.ID(), time.Now())
	case *gonet.SessionClose:
		logs.Warn("session_%d close at %v", msg.Session.ID(), time.Now())
	case *proto.Ping:
		logs.Info("session_%d ping at %d", msg.Session.ID(), data.At)
	}
}
```
## 在线游戏demo
- **使用etcd+mysql+beego+gonet+cocos creator制作**  
服务端：大厅服+游戏服+服务注册  
客户端：大厅+子游戏模式  
http://116.62.245.150:8087/web-desktop/
![display](./display_lkby.gif)
## 测试
## FAQ
## 参与贡献
#### QQ群：795611332

