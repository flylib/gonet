package demo

// NOTE: The quic transport (transport/quic) session type is unexported.
// These tests are commented out pending export of the session type or
// provision of a factory function from the quic package.
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
// 	transport "github.com/flylib/gonet/transport/quic"
// 	"github.com/flylib/goutils/codec/json"
// 	"github.com/flylib/pkg/log/builtinlog"
// 	"log"
// 	"testing"
// 	"time"
// )

// func TestQuicServer(t *testing.T) { ... }
// func TestQuicClient(t *testing.T) { ... }
