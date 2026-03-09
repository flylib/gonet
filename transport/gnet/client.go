package gnet

import (
	"github.com/flylib/gonet"
	"github.com/panjf2000/gnet/v2"
	"log"
	"time"
)

type client struct {
	gonet.PeerCommon[*session]
	eng gnet.Engine
	cli *gnet.Client
}

func NewClient(ctx *gonet.Context[*session], options ...Option) gonet.IClient {
	var opt option
	for _, f := range options {
		f(&opt)
	}
	opt.Logger = ctx.GetLogger()

	c := &client{}
	cli, err := gnet.NewClient(c, gnet.WithOptions(opt.Options))
	if err != nil {
		panic(err)
	}
	c.cli = cli
	c.WithContext(ctx)
	return c
}

func (c *client) OnBoot(eng gnet.Engine) (action gnet.Action) {
	c.eng = eng
	return gnet.None
}

func (c *client) OnShutdown(eng gnet.Engine) {
}

func (c *client) OnOpen(conn gnet.Conn) (out []byte, action gnet.Action) {
	panic("implement me")
}

func (c *client) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	is, ok := c.GetCtx().GetSession(uint64(conn.Fd()))
	if ok {
		c.GetCtx().RecycleSession(is)
	}
	return gnet.None
}

func (c *client) OnTraffic(conn gnet.Conn) (action gnet.Action) {
	buf, err := conn.Next(-1)
	if err != nil {
		return gnet.Close
	}
	is, ok := c.GetCtx().GetSession(uint64(conn.Fd()))
	if !ok {
		return gnet.Close
	}
	message, _, err := c.GetCtx().UnPackage(is, buf)
	if err != nil {
		log.Printf("session_%v msg parser error,reason is %v \n", conn.Fd(), err)
		return gnet.None
	}
	c.GetCtx().PushGlobalMessageQueue(message)
	return gnet.None
}

func (c *client) OnTick() (delay time.Duration, action gnet.Action) {
	return time.Hour, gnet.None
}

func (c *client) Dial(addr string) (gonet.ISession, error) {
	conn, err := c.cli.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return newSession(c.GetCtx(), conn), nil
}

func (c *client) Close() error {
	return c.cli.Stop()
}
