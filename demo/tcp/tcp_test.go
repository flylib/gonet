package tcp

import (
	"demo/handler"
	"demo/proto"
	"fmt"
	"github.com/flylib/gonet"
	transport "github.com/flylib/gonet/transport/tcp"
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/pkg/log/builtinlog"
	"log"
	"testing"
	"time"
)

var addr = "localhost:8089"

func TestTcpServer(t *testing.T) {
	gonet.SetupContext(
		gonet.WithEventHandler(handler.EventHandler{}),
		gonet.WithNetPackager(gonet.TcpNetPackager{}),

		gonet.MustWithSessionType(transport.SessionType()),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	if err := transport.NewServer().Listen(addr); err != nil {
		log.Fatal(err)
	}
}

func TestTcpClient(t *testing.T) {
	gonet.SetupContext(
		gonet.WithEventHandler(handler.EventHandler{}),
		gonet.WithNetPackager(gonet.TcpNetPackager{}),

		gonet.MustWithSessionType(transport.SessionType()),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	session, err := transport.NewClient(transport.WithHandshakeTimeout(5 * time.Second)).Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connect success")

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
