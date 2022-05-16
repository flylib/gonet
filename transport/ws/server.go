package ws

import (
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"net/http"
	"net/url"
)

//接收端
type server struct {
	ServerIdentify
	//指定将HTTP连接升级到WebSocket连接的参数。
	upGrader websocket.Upgrader
	//响应头
	//respHeader http.Header
}

func (s *server) Start() error {
	url, err := url.Parse(s.Addr())
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc(url.Path, s.newConn)
	return http.ListenAndServe(url.Host, mux)
}

func (s *server) Stop() error {
	// TODO 关闭处理
	return nil
}

func (s *server) newConn(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Connection", "upgrade") //升级
	r.Header.Add("Upgrade", "websocket")  //websocket
	conn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	newConnection := newConn(conn)
	go newConnection.recvLoop()
}
