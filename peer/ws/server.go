package ws

import (
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	. "github.com/zjllib/goNet"
	"net/http"
	"net/url"
)

//接收端
type server struct {
	PeerIdentify
	//指定将HTTP连接升级到WebSocket连接的参数。
	upGrader websocket.Upgrader
	//响应头
	//respHeader http.Header
}

func init() {
	identify := PeerIdentify{}
	identify.SetType(PeertypeServer)
	//响应头
	//var header http.Header = make(map[string][]string)
	//header.Add("Access-Control-Allow-Origin", "*")
	s := &server{
		PeerIdentify: identify,
		upGrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	RegisterPeer(s)
}

func (s *server) Start() {
	url, err := url.Parse(s.Addr())
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc(url.Path, s.newConn)
	logs.Info("#websocket.listen(%s)", s.Addr())

	err = http.ListenAndServe(url.Host, mux)
	if err != nil {
		panic(err)
	}
}

func (s *server) newConn(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Connection", "upgrade") //升级
	r.Header.Add("Upgrade", "websocket")  //websocket
	conn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error("http covert to websocket err:", err.Error())
		return
	}
	logs.Info("new connect from ", conn.RemoteAddr())
	session := newSession(conn)

	go session.recvLoop()
	go session.sendLoop()
}

func (s *server) Stop() {
	// TODO 关闭处理
}
