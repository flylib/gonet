package demo

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/flylib/gonet"
	transport "github.com/flylib/gonet/transport/fastws"
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/pkg/log/builtinlog"
)

// startFastwsBenchServer starts a fastws server on a random port.
func startFastwsBenchServer(handler gonet.IEventHandler) (string, func()) {
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

// dialFastwsBenchClient creates a fastws client and dials the given address.
func dialFastwsBenchClient(addr string, handler gonet.IEventHandler) gonet.ISession {
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

// BenchmarkFastwsSendRecv measures single-goroutine send→recv round-trip throughput.
func BenchmarkFastwsSendRecv(b *testing.B) {
	var received atomic.Int64

	serverHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {
		msg.From().Send(benchMsgID, &benchMsg{Data: benchPayload})
	}}
	addr, cleanup := startFastwsBenchServer(serverHandler)
	defer cleanup()

	clientCh := make(chan struct{}, 1)
	clientHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {
		received.Add(1)
		select {
		case clientCh <- struct{}{}:
		default:
		}
	}}
	session := dialFastwsBenchClient(addr, clientHandler)

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

// BenchmarkFastwsThroughput measures one-way send throughput.
func BenchmarkFastwsThroughput(b *testing.B) {
	var received atomic.Int64
	var wg sync.WaitGroup
	wg.Add(b.N)

	serverHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {
		received.Add(1)
		wg.Done()
	}}
	addr, cleanup := startFastwsBenchServer(serverHandler)
	defer cleanup()

	clientHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {}}
	session := dialFastwsBenchClient(addr, clientHandler)

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

// BenchmarkFastwsParallelSend measures parallel send throughput from multiple goroutines.
func BenchmarkFastwsParallelSend(b *testing.B) {
	var received atomic.Int64
	var wg sync.WaitGroup
	wg.Add(b.N)

	serverHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {
		received.Add(1)
		wg.Done()
	}}
	addr, cleanup := startFastwsBenchServer(serverHandler)
	defer cleanup()

	clientHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {}}
	session := dialFastwsBenchClient(addr, clientHandler)

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

// TestFastwsBenchmarkReport runs a quick throughput test and prints a summary.
func TestFastwsBenchmarkReport(t *testing.T) {
	const total = 50000

	var received atomic.Int64
	var wg sync.WaitGroup
	wg.Add(total)

	serverHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {
		received.Add(1)
		wg.Done()
	}}
	addr, cleanup := startFastwsBenchServer(serverHandler)
	defer cleanup()

	clientHandler := &benchHandler{onMsg: func(msg gonet.IMessage) {}}
	session := dialFastwsBenchClient(addr, clientHandler)

	start := time.Now()
	for i := 0; i < total; i++ {
		if err := session.Send(benchMsgID, &benchMsg{Data: benchPayload}); err != nil {
			t.Fatal(err)
		}
	}
	wg.Wait()
	elapsed := time.Since(start)

	qps := float64(total) / elapsed.Seconds()
	fmt.Printf("\n=== FastWS Throughput Report ===\n")
	fmt.Printf("  Messages:  %d\n", total)
	fmt.Printf("  Elapsed:   %v\n", elapsed)
	fmt.Printf("  QPS:       %.0f msg/s\n", qps)
	fmt.Printf("  Avg:       %v/msg\n", elapsed/time.Duration(total))
	fmt.Printf("================================\n\n")
}
