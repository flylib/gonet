
![goNetlogo](./logo.jpg)
## version
 v 1.0.0
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
git clone https://github.com/zjllib/goNet.git
```

## 使用样例参考
- 服务端
```go
	p := goNet.NewPeer(
		goNet.Options{
			Addr:          "ws://:8083/echo",
			PeerType:      goNet.PEERTYPE_SERVER,
		})
	p.Start()
```
- 客户端
```go
	p := goNet.NewPeer(
				goNet.Options{
					Addr: "ws://:8083/echo",
					PeerType: goNet.PEERTYPE_CLIENT,
				})
			p.Start()
```
- 消息处理实现及注册
```go
//101~200  登录注册模块
const (
	MsgID_LoginReq = 101
	MsgID_LoginOut = 102
)

//消息注入
func init() {
	goNet.RegisterMsg(MsgID_LoginReq, goNet.SYSTEM_CONTROLLER_IDX, LoginReq{})
	goNet.RegisterMsg(MsgID_LoginOut, goNet.SYSTEM_CONTROLLER_IDX, LoginOut{})
}
////实现消息控制器
//type Controller interface {
//  	OnMsg(session Session, msg interface{})
 // }

func (u *Controller) OnMsg(session goNet.Session, data interface{}) {
	switch msg := data.(type) {
	case *proto.LoginReq: //登录请求
	  //todo something
	case *proto.LoginOut: //退出登录
    //todo something
	}
}
```
## 在线游戏demo
- **使用etcd+mysql+beego+goNet+cocos creator制作**  
服务端：大厅服+游戏服+服务注册  
客户端：大厅+子游戏模式  
http://116.62.245.150:8087/web-desktop/
![display](./display_lkby.gif)
## 测试
## FAQ
## 参与贡献
#### QQ群：795611332

