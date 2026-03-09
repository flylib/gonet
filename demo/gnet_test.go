package demo

// NOTE: The gnet transport (transport/gnet) has not been updated to the generic
// Context[S] API yet. The tests below are commented out until that transport is
// refactored to use gonet.Context[*gnet.Session] and gonet.PeerCommon[*gnet.Session].
//
// When the gnet transport is updated, the usage pattern will be:
//
//   ctx := gonet.NewContext(
//       func() *transport.Session { return new(transport.Session) },
//       gonet.WithEventHandler(handler.EventHandler{}),
//       gonet.MustWithCodec(&json.Codec{}),
//       gonet.MustWithLogger(builtinlog.NewLogger()),
//   )
//   transport.NewServer(ctx).Listen("localhost:8088")

// import (
// 	"demo/handler"
// 	"demo/proto"
// 	"fmt"
// 	"github.com/flylib/gonet"
// 	transport "github.com/flylib/gonet/transport/gnet"
// 	"github.com/flylib/goutils/codec/json"
// 	"github.com/flylib/pkg/log/builtinlog"
// 	"log"
// 	"testing"
// 	"time"
// )

// func TestGNETServer(t *testing.T) {
// 	ctx := gonet.NewContext(
// 		func() *transport.Session { return new(transport.Session) },
// 		gonet.WithEventHandler(handler.EventHandler{}),
// 		gonet.MustWithCodec(&json.Codec{}),
// 		gonet.MustWithLogger(builtinlog.NewLogger()),
// 	)
// 	t.Log("server listen on localhost:8088")
// 	if err := transport.NewServer(ctx).Listen("localhost:8088"); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func TestGNETClient(t *testing.T) {
// 	ctx := gonet.NewContext(
// 		func() *transport.Session { return new(transport.Session) },
// 		gonet.WithEventHandler(handler.EventHandler{}),
// 		gonet.MustWithCodec(&json.Codec{}),
// 		gonet.MustWithLogger(builtinlog.NewLogger()),
// 	)
// 	session, err := transport.NewClient(ctx).Dial("localhost:8088")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	t.Log("connect success")
//
// 	tick := time.Tick(time.Second * 1)
// 	var i int
// 	for range tick {
// 		i++
// 		err = session.Send(101, &proto.Say{
// 			fmt.Sprintf("hello server %d", i),
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }
