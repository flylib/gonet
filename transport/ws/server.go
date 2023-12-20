package ws

import (
	"github.com/flylib/gonet"
	"net/http"
	"net/url"
	"time"
)

var _ gonet.IServer = new(server)

// 接收端
type server struct {
	gonet.PeerCommon
	option
}

func NewServer(ctx *gonet.Context, options ...Option) gonet.IServer {
	s := &server{}
	for _, f := range options {
		f(&s.option)
	}
	s.WithContext(ctx)
	return s
}

func (s *server) Listen(addr string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}
	s.SetAddr(addr)
	mux := http.NewServeMux()
	mux.HandleFunc(u.Path, s.newConn)
	return http.ListenAndServe(u.Host, mux)
}

func (s *server) Close() error {
	s.option.Upgrader.HandshakeTimeout = time.Nanosecond
	return nil
}

func (s *server) newConn(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Connection", "upgrade") //升级
	r.Header.Add("Upgrade", "websocket")  //websocket

	conn, err := s.option.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	ns := newSession(s.Context, conn)
	go ns.ReadLoop()
}
