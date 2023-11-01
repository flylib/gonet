
![gonetlogo](docs/logo.jpg)
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
		gonet.WorkerPoolMaxSize(20),
		handler.InitServerRouter,
	)

	if err := ws.NewServer(ctx).Listen("ws://localhost:8088/center/ws"); err != nil {
		log.Fatal(err)
	}
}

// handler.go package
// 消息路由
func InitServerRouter(ctx *gonet.Context) error {
	ctx.Route(gonet.SessionConnect, nil, serverHandler)
	ctx.Route(gonet.SessionClose, nil, serverHandler)
	ctx.Route(101, proto.Say{}, serverHandler)
	return nil
}

func serverHandler(s gonet.ISession, msg gonet.IMessage) {
	switch msg.ID() {
	case gonet.SessionConnect:
		log.Println("connected session_id:", s.ID(), " ip:", s.RemoteAddr().String())
	case gonet.SessionClose:
		log.Println("connected session_id:", s.ID(), " error:", msg.Body())
	case 101:
		fmt.Println("session_id:", s.ID(), " say ", msg.Body().(*proto.Say).Content)
		err := s.Send(proto.Say{Content: "hell client"})
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Println("unknown message id:", msg.ID())
	}
}
```


## 测试
## FAQ
## 参与贡献
#### QQ群：795611332

