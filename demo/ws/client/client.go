package main

import (
	"github.com/Quantumoffices/goNet"
	_ "github.com/Quantumoffices/goNet/codec/json"
	_ "github.com/Quantumoffices/goNet/peer/ws"
	"time"
)

//ws://47.57.65.221:8088/game/blockInfo
//ws://192.168.0.125:8088/game/blockInfo
func main() {
	p := goNet.NewPeer(
		goNet.Options{
			Addr:     "ws://192.168.0.125:2020/center/ws",
			PeerType: goNet.PEERTYPE_CLIENT,
			//ReadDeadline:  0,
			//WriteDeadline: 0,
			//PoolSize:      0,
			//PanicHandler:  nil,
			//AllowMaxConn:  0,
		})
	p.Start()
	s := goNet.SessionManager.GetSessionById(1)
	for {
		time.Sleep(time.Second)
		s.Send(goNet.Ping{TimeStamp: time.Now().Unix()})
	}
}
