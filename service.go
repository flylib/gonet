package gonet

import "github.com/zjllib/gonet/v3/transport"

/*
+---------------------------------------------------+
+				     service						+
+---------------------------------------------------+
+		server			|		client				+
+---------------------------------------------------+
+		bee worker、conn pool、codec					+
+---------------------------------------------------+
+		transport(udp、tcp、ws、quic)			    +
+---------------------------------------------------+
*/

//一切皆服务
type IService interface {
	//服务名
	Name() string
	// 开启服务
	Start() error
	// 停止服务
	Stop() error
	// Client is used to call services
	Client() transport.IClient
	// Server is for handling requests and events
	Server() transport.IServer
}
