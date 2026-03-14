package demo

import (
	"demo/handler"
	"demo/proto"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/flylib/gonet"
	transport "github.com/flylib/gonet/transport/fastws"
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/pkg/log/builtinlog"
)

func TestFastwsServer(t *testing.T) {
	ctx := gonet.NewAppContext(
		func() *transport.Session { return new(transport.Session) },
		gonet.WithEventHandler(handler.EventHandler{}),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	fmt.Println("server listen on ws://localhost:8089/ws")
	if err := transport.NewServer(ctx).Listen("ws://localhost:8089/ws"); err != nil {
		log.Fatal(err)
	}
}

func TestFastwsClient(t *testing.T) {
	ctx := gonet.NewAppContext(
		func() *transport.Session { return new(transport.Session) },
		gonet.WithEventHandler(handler.EventHandler{}),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	session, err := transport.NewClient(ctx, transport.WithHandshakeTimeout(5*time.Second)).Dial("ws://localhost:8089/ws")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connect success")

	tick := time.Tick(time.Second * 1)
	var i int
	for range tick {
		i++
		err = session.Send(101, &proto.Say{
			Content: fmt.Sprintf("hello server %d", i),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
