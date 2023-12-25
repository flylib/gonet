package demo

import (
	"demo/handler"
	"demo/proto"
	"fmt"
	"github.com/flylib/gonet"
	transport "github.com/flylib/gonet/transport/quic"
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/pkg/log/builtinlog"
	"log"
	"testing"
	"time"
)

func TestQuicServer(t *testing.T) {
	ctx := gonet.NewContext(
		gonet.WithEventHandler(handler.EventHandler{}),

		gonet.MustWithSessionType(transport.SessionType()),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	t.Log("server listen on localhost:8088")
	if err := transport.NewServer(ctx).Listen("localhost:8088"); err != nil {
		log.Fatal(err)
	}
}

func TestQuicClient(t *testing.T) {
	ctx := gonet.NewContext(
		gonet.WithEventHandler(handler.EventHandler{}),

		gonet.MustWithSessionType(transport.SessionType()),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	session, err := transport.NewClient(ctx).Dial("localhost:8088")
	if err != nil {
		log.Fatal(err)
	}

	t.Log("connect success")

	tick := time.Tick(time.Second * 1)
	var i int
	for range tick {
		//fmt.Println("send msg", i)
		i++
		err = session.Send(101, &proto.Say{
			fmt.Sprintf("hello server %d", i),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
