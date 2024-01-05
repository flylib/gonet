package kcp

import (
	"crypto/sha1"
	"github.com/flylib/gonet"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

type client struct {
	gonet.PeerCommon
	conn *kcp.UDPSession
	option
}

func NewClient(options ...Option) gonet.IClient {
	c := &client{}
	for _, f := range options {
		f(&c.option)
	}
	return c
}

func (c *client) Dial(addr string) (gonet.ISession, error) {
	key := pbkdf2.Key([]byte(c.PBKDF2Password), []byte(c.PBKDF2Salt), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	conn, err := kcp.DialWithOptions(addr, block, 10, 3)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	s := newSession(conn)
	go s.readLoop()
	return s, err
}

func (c *client) Close() error {
	return c.conn.Close()
}
