
## version
 `v3.1.0`
 
## 介绍
一个基于go语言开发的网络脚手架,参考[cellnet](https://github.com/davyxu/cellnet)和[gin](https://github.com/gin-gonic/gin) 两大开框架的设计，使用非常方便简洁，轻松让你开发出高并发高性能的网络应用，可以用于游戏,app等任何领域的通讯。

## 主要技术理念
- Session pool 会话池
- Goroutine pool  消息处理协程池
- Message cache layer 消息缓存层
- Message sequence 消息排序

## 架构图
![architecture](./architecture.png)


## 主要特性及追求目标
- 高并发
- 高性能
- 简单易用
- 线性安全
- 兼容性强
- 高度可配置
- 多领域应用
- 防崩溃
- 错误快速定位

## 通讯协议支持
- [x] TCP
- [x] UDP
- [x] WEBSOCKET
- [x] QUIC
- [x] KCP

## 数据编码格式支持
- [x] json
- [x] xml
- [x] binary
- [x] protobuf


## 安装教程
### **1.** git clone到 GOPATH/src目录下

```
git clone https://github.com/flylib/gonet.git
```

## 使用样例参考
```go
//main.go
package main

import (
	"github.com/flylib/gonet"
	"github.com/flylib/gonet/demo/handler"
	"github.com/flylib/gonet/transport/ws" //协议
	"log"
)

func main() {
	ctx := gonet.NewContext(
		gonet.WithEventHandler(handler.EventHandler{}),

		gonet.MustWithSessionType(transport.SessionType()),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	fmt.Println("server listen on ws://localhost:8088/center/ws")
	if err := transport.NewServer(ctx).Listen("ws://localhost:8088/center/ws"); err != nil {
		log.Fatal(err)
	}
}


type EventHandler struct {
}

func (e EventHandler) OnError(session gonet.ISession, err error) {
	println(fmt.Sprintf("sesson-%d error-%v", session, err))
}

func (e EventHandler) OnConnect(session gonet.ISession) {
	fmt.Println(fmt.Sprintf("new session-%d from-%s", session.ID(), session.RemoteAddr()))
}

func (e EventHandler) OnClose(session gonet.ISession, err error) {
	//TODO implement me
	fmt.Println(fmt.Sprintf("session close-%d", session.ID()))
}

func (e EventHandler) OnMessage(message gonet.IMessage) {
	//TODO implement me
	switch message.ID() {
	case 101:
		fmt.Println("session-", message.From().ID(), " say:", string(message.Body()))
	}
}

```


## 测试
## FAQ
## 参与贡献
#### QQ群：795611332

