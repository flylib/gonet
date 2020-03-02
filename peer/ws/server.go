package ws

import (
	"github.com/gorilla/websocket"
	"goNet"
	. "goNet/log"
	"net/http"
	"net/url"
)

//接收端
type server struct {
	goNet.PeerIdentify
	certfile string
	keyfile  string
	//指定将HTTP连接升级到WebSocket连接的参数。
	upgrader websocket.Upgrader
	//响应头
	//respHeader http.Header
}

func init() {
	identify := goNet.PeerIdentify{}
	identify.SetType(goNet.PEER_SERVER)

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
	goNet.RegisterPeer(s)
}

//wss加密通信协议
func (s *server) SetHttps(certfile, keyfile string) {
	s.certfile = certfile
	s.keyfile = keyfile
}

func (s *server) Start() {
	url, err := url.Parse(s.Addr())
	if err != nil {
		Log.Fatalf("#websocket.url parse failed(%s) %v", s.Addr(), err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc(url.Path, s.newConn)
	Log.Infof("#websocket.listen(%s)", s.Addr())

	if url.Scheme == "https" {
		err = http.ListenAndServeTLS(url.Host, s.certfile, s.keyfile, mux)
	} else {
		err = http.ListenAndServe(url.Host, mux)
	}
	if err != nil {
		Log.Fatalf("#websocket stop listen , failed(%s) %v", s.Addr(), err.Error())
	}
}

func (s *server) newConn(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		Log.Error("http covert to websocket err:", err.Error())
		return
	}
	Log.Info("new connect from ", conn.RemoteAddr())
	ses := newSession(conn)
	ses.recvLoop()
}

func (s *server) Stop() {
	// TODO 关闭处理
}
