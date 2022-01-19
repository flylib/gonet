
![gonetlogo](docs/logo.jpg)
## version
 v3
 
## 介绍
一个基于go语言开发的网络脚手架,参考[cellnet](https://github.com/davyxu/cellnet)和[gin](https://github.com/gin-gonic/gin) 两大开框架的设计，使用非常方便简洁，轻松让你开发出高并发高性能的网络应用，可以用于游戏,app等任何领域的通讯。

## 主要技术理念
- Session pool 会话池
- Routine pool  消息处理协程池
- Message cache layer 消息缓存层
- Message ID route 依据消息ID进行路由


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
```go
	package main
    
    import (
    	"fmt"
    	"github.com/zjllib/gonet/v3"
    	"github.com/zjllib/gonet/v3/demo/proto"
    	_ "github.com/zjllib/gonet/v3/transport/ws" //协议
    	"log"
    )
    //消息路由
    func init() {
    	gonet.Route(gonet.SessionConnect, nil, Handler)
    	gonet.Route(gonet.SessionClose, nil, Handler)
    	gonet.Route(101, proto.Say{}, Handler)
    }
    
    func main() {
    	server := gonet.NewServer(
    		gonet.Address("ws://localhost:8088/center/ws"), //listen addr
    		gonet.MaxWorkerPoolSize(20))
    	log.Printf("server listening on %s", server.Addr())
    	if err := server.Start(); err != nil {
    		log.Fatal(err)
    	}
    }
    //消息处理函数
    func Handler(msg *gonet.Message) {
    	switch msg.ID {
    	case gonet.SessionConnect:
    		log.Println("connected session_id:", msg.Session.ID(), " ip:", msg.Session.RemoteAddr().String())
    	case gonet.SessionClose:
    		log.Println("connected session_id:", msg.Session.ID(), " error:", msg.Body)
    	case 101:
    		fmt.Println("session_id:", msg.Session.ID(), " say ", msg.Body.(*proto.Say).Content)
    		//fmt.Println(reflect.TypeOf(msg.Body))
    	default:
    		log.Println("unknown session_id:", msg.ID)
    	}
    }

```


## 测试
## FAQ
## 参与贡献
#### QQ群：795611332

