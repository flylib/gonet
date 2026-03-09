package demo

// NOTE: The udp transport (transport/udp) session type is unexported.
// These tests are commented out pending export of the session type or
// provision of a factory function from the udp package.
//
// When ready, the usage pattern will be:
//
//   ctx := gonet.NewContext(
//       transport.NewSession,
//       gonet.WithEventHandler(handler.EventHandler{}),
//       gonet.MustWithCodec(&json.Codec{}),
//       gonet.MustWithLogger(builtinlog.NewLogger()),
//   )

// import (
// 	"demo/handler"
// 	"demo/proto"
// 	"fmt"
// 	"github.com/flylib/gonet"
// 	transport "github.com/flylib/gonet/transport/udp"
// 	"github.com/flylib/goutils/codec/json"
// 	"github.com/flylib/pkg/log/builtinlog"
// 	"log"
// 	"testing"
// 	"time"
// )

// func TestUDPServer(t *testing.T) { ... }
// func TestUDPClient(t *testing.T) { ... }
