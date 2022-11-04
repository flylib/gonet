package ws

import (
	"github.com/gorilla/websocket"
	"github.com/zjllib/gonet/v3/transport"
	"net/http"
	"net/url"
	"reflect"
)

var _ transport.Transport = new(ws)

//接收端
type ws struct {
	transport.TransportIdentify
	//指定将HTTP连接升级到WebSocket连接的参数。
	upGrader websocket.Upgrader
	//响应头
	//respHeader http.Header
}

func NewTransport(addr string) *ws {
	t := &ws{}
	t.SetAddr(addr)
	return t
}

func (s *ws) Listen() error {
	url, err := url.Parse(s.Addr())
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc(url.Path, s.newConn)
	return http.ListenAndServe(url.Host, mux)
}

func (s *ws) Stop() error {
	// TODO 关闭处理
	return nil
}

func (s *ws) SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}

func (s *ws) newConn(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Connection", "upgrade") //升级
	r.Header.Add("Upgrade", "websocket")  //websocket

	conn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go newSession(conn).recvLoop()
}
