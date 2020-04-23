package ws

import (
	"github.com/Quantumoffices/beego/logs"
	. "github.com/Quantumoffices/goNet"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

//接收端
type server struct {
	PeerIdentify
	certfile string
	keyfile  string
	//指定将HTTP连接升级到WebSocket连接的参数。
	upgrader websocket.Upgrader
	//响应头
	//respHeader http.Header
}

func init() {
	identify := PeerIdentify{}
	identify.SetType(PEERTYPE_SERVER)
	////响应头
	//var header http.Header = make(map[string][]string)
	//header.Add("Access-Control-Allow-Origin", "*")

	s := &server{
		PeerIdentify: identify,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		//respHeader: header,
	}
	RegisterPeer(s)
}

//wss加密通信协议
func (s *server) SetHttps(certfile, keyfile string) {
	s.certfile = certfile
	s.keyfile = keyfile
}

func (s *server) Start() {
	url, err := url.Parse(s.Addr())
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc(url.Path, s.newConn)

	if url.Scheme == "https" {
		err = http.ListenAndServeTLS(url.Host, s.certfile, s.keyfile, mux)
	} else {
		err = http.ListenAndServe(url.Host, mux)
	}
	if err != nil {
		panic(err)
	}
}

func (s *server) newConn(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error("http covert to websocket err:", err.Error())
		return
	}
	logs.Info("new connect from ", conn.RemoteAddr())
	ses := newSession(conn)
	ses.recvLoop()
}

func (s *server) Stop() {
	// TODO 关闭处理
}
