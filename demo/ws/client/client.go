package main

import (
	"github.com/zjllib/goNet"
	_ "github.com/zjllib/goNet/codec/json"
	_ "github.com/zjllib/goNet/peer/ws"
)

//ws://47.57.65.221:8088/game/blockInfo
//ws://192.168.0.125:8088/game/blockInfo

// /lottery/api/v1/ws
func main() {

	client := goNet.NewClient("ws://localhost:8088/center/ws")
	client.Start()
	//for {
	//	//client := goNet.NewClient("ws://localhost:8088/center/ws")
	//	//client.Start()
	//		session, ok := goNet.FindSession(uint64(goNet.SessionCount()))
	//		if ok {
	//			for {
	//				time.Sleep(time.Second)
	//				session.Send(goNet.Ping{TimeStamp: time.Now().Unix()})
	//				logs.Info("send---- ")
	//			}
	//		}
	//	}()
	//	time.Sleep(time.Second)
	//}

	//p := goNet.NewPeer(
	//	goNet.Options{
	//		Addr: "ws://192.168.0.125:8083/center/ws",
	//		//Addr:     "ws://192.168.0.125:4160/lottery/api/v1/ws",
	//		PeerType: goNet.PEERTYPE_CLIENT,
	//		//ReadDeadline:  0,
	//		//WriteDeadline: 0,
	//		//PoolSize:      0,
	//		//PanicHandler:  nil,
	//		//AllowMaxConn:  0,
	//	})
	//p.Start()
	//session, ok := goNet.FindSession(uint64(goNet.SessionCount()))
	//if ok {
	//	for {
	//		time.Sleep(time.Second)
	//		session.Send(goNet.Ping{TimeStamp: time.Now().Unix()})
	//		logs.Info("send---- ")
	//	}
	//}
}
