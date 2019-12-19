
![goNetlogo](./logo.jpg)
## version
 v 1.0.0
## 介绍
一个基于go语言开发的网络脚手架,参考[cellnet](https://github.com/davyxu/cellnet)和[beego](https://github.com/astaxie/beego)两大开框架的设计，轻松上手，轻松让你开发出高并发高性能的网络应用，可以用于游戏,app等任何领域的通讯。

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
git clone https://github.com/Quantumoffices/goNet.git
```

## 使用样例参考
- 服务端
```go
package main

import (
	"goNet"
	_ "goNet/codec/json"
	_ "goNet/peer/tcp"
)

func main() {
	p := goNet.NewPeer("server",":8087")
	p.Start()
}

```
- 客户端
```go
package main

import (
	"goNet"
	_ "goNet/codec/json"
	_ "goNet/peer/tcp"
)

func main() {
	p := goNet.NewPeer("client", ":8087")
	p.Start()
       //todo something
}
```
- 消息处理实现及注册
```go
package msg
import (
	"goNet"
)
//消息注册
func init() {
	goNet.RegisterMessage(0, Ping{})
	goNet.RegisterMessage(1, Pong{})
}

//心跳
type Ping struct {
	TimeStamp int64 `json:"time_stamp",xml:"time_stamp"`
}
type Pong struct {
	TimeStamp int64 `json:"time_stamp",xml:"time_stamp"`
}
type SessionClose struct {
}

//消息需要实现 goNet.Message接口
func (p *Ping) Handle(session Session) {
	logrus.Infof("session_%v ping at time=%v", session.ID(), time.Unix(p.TimeStamp, 0).String())
	session.Send(Pong{TimeStamp: time.Now().Unix(),})
}
func (p *Pong) Handle(session Session) {
	logrus.Infof("session_%v pong at time=%v", session.ID(), time.Unix(p.TimeStamp, 0).String())
}

func (s *SessionClose) Handle(session Session) {
	logrus.Errorf("session_%v closed", session.ID())
}
```
## 在线游戏demo
- **使用etcd+mysql+beego+goNet+cocos creator制作**  
服务端：大厅服+游戏服+服务注册  
客户端：大厅+子游戏模式  
http://www.quantumstudio:8080/client-release-signed.apk
![display](./display_lkby.gif)
## 测试
## FAQ
## 参与贡献
## 感谢您的帮助
我的母亲今年(2019年9月份)被查出患肺腺癌晚期，目前家里入不敷出，如果框架能帮助到您，希望帮助一下我的母亲！在此真诚的谢谢您并在交流群@我,会赠与游戏demo源码！  
<img src="./ali_pay.JPG" height="300">  <img src="./wx_pay.JPG" height="300">


#### QQ交流群：795611332

