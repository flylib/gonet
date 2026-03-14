package demo

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/flylib/gonet"
	transport "github.com/flylib/gonet/transport/gorillaws"
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/pkg/log/builtinlog"
)

const (
	benchMsgID   = 1
	benchPayload = "hello"
)

type benchMsg struct {
	Data string `json:"data"`
}

// benchHandler is a minimal event handler for benchmarking.
type benchHandler struct {
	onMsg func(gonet.IMessage)
}

func (h *benchHandler) OnConnect(gonet.ISession)      {}
func (h *benchHandler) OnClose(gonet.ISession, error) {}
func (h *benchHandler) OnError(gonet.ISession, error) {}
func (h *benchHandler) OnMessage(msg gonet.IMessage)  { h.onMsg(msg) }

// freePort picks an available TCP port from the OS.
func freePort() int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

// startBenchServer starts a WS server on a random port, returns the address and cleanup func.
func startBenchServer(handler gonet.IEventHandler) (string, func()) {
	ctx := gonet.NewAppContext(
		func() *transport.Session { return new(transport.Session) },
		gonet.WithEventHandler(handler),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	addr := fmt.Sprintf("ws://127.0.0.1:%d/ws", freePort())
	srv := transport.NewServer(ctx)
	go srv.Listen(addr)
	time.Sleep(200 * time.Millisecond)
	return addr, func() { srv.Close() }
}

// dialBenchClient creates a client and dials the given address, returning the session.
func dialBenchClient(addr string, handler gonet.IEventHandler) gonet.ISession {
	ctx := gonet.NewAppContext(
		func() *transport.Session { return new(transport.Session) },
		gonet.WithEventHandler(handler),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	session, err := transport.NewClient(ctx, transport.WithHandshakeTimeout(5*time.Second)).Dial(addr)
	if err != nil {
		panic(err)
	}
	return session
}

// BenchmarkWsSendRecv measures single-goroutine send→recv round-trip throughput.
func BenchmarkWsSendRecv(b *testing.B) {
	var received atomic.Int64

	serverHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {
		// echo back
		msg.From().Send(benchMsgID, &benchMsg{Data: benchPayload})
	}}
	addr, cleanup := startBenchServer(serverHandler)
	defer cleanup()

	clientCh := make(chan struct{}, 1)
	clientHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {
		received.Add(1)
		select {
		case clientCh <- struct{}{}:
		default:
		}
	}}
	session := dialBenchClient(addr, clientHandler)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if err := session.Send(benchMsgID, &benchMsg{Data: benchPayload}); err != nil {
			b.Fatal(err)
		}
		<-clientCh
	}

	b.StopTimer()
	b.Logf("total messages echoed: %d", received.Load())
}

// BenchmarkWsThroughput measures one-way send throughput (fire-and-forget).
func BenchmarkWsThroughput(b *testing.B) {
	var received atomic.Int64
	var wg sync.WaitGroup
	wg.Add(b.N)

	serverHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {
		received.Add(1)
		wg.Done()
	}}
	addr, cleanup := startBenchServer(serverHandler)
	defer cleanup()

	clientHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {}}
	session := dialBenchClient(addr, clientHandler)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if err := session.Send(benchMsgID, &benchMsg{Data: benchPayload}); err != nil {
			b.Fatal(err)
		}
	}

	wg.Wait()
	b.StopTimer()
	b.Logf("total messages received: %d", received.Load())
}

// BenchmarkWsParallelSend measures parallel send throughput from multiple goroutines.
func BenchmarkWsParallelSend(b *testing.B) {
	var received atomic.Int64
	var wg sync.WaitGroup
	wg.Add(b.N)

	serverHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {
		received.Add(1)
		wg.Done()
	}}
	addr, cleanup := startBenchServer(serverHandler)
	defer cleanup()

	clientHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {}}
	session := dialBenchClient(addr, clientHandler)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := session.Send(benchMsgID, &benchMsg{Data: benchPayload}); err != nil {
				b.Error(err)
				return
			}
		}
	})

	wg.Wait()
	b.StopTimer()
	b.Logf("total messages received: %d", received.Load())
}

// TestWsBenchmarkReport runs a quick throughput test and prints a summary.
func TestWsBenchmarkReport(t *testing.T) {
	const total = 50000

	var received atomic.Int64
	var wg sync.WaitGroup
	wg.Add(total)

	serverHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {
		received.Add(1)
		wg.Done()
	}}
	addr, cleanup := startBenchServer(serverHandler)
	defer cleanup()

	clientHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {}}
	session := dialBenchClient(addr, clientHandler)

	start := time.Now()
	for i := 0; i < total; i++ {
		if err := session.Send(benchMsgID, &benchMsg{Data: benchPayload}); err != nil {
			t.Fatal(err)
		}
	}
	wg.Wait()
	elapsed := time.Since(start)

	qps := float64(total) / elapsed.Seconds()
	fmt.Printf("\n=== WebSocket Throughput Report ===\n")
	fmt.Printf("  Messages:  %d\n", total)
	fmt.Printf("  Elapsed:   %v\n", elapsed)
	fmt.Printf("  QPS:       %.0f msg/s\n", qps)
	fmt.Printf("  Avg:       %v/msg\n", elapsed/time.Duration(total))
	fmt.Printf("==================================\n\n")
}
