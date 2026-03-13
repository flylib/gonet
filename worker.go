package gonet

import (
	"runtime"
	"runtime/debug"

	ilog "github.com/flylib/interface/log"
)

type poolConfig struct {
	queueSize  int32
	maxNum     int32 // reserved, unused in shard model
	maxIdleNum int32 // shard count (= worker count); default: NumCPU
}

// GoroutinePool is a sharded goroutine pool for processing messages.
//
// Messages are routed to a shard by (sessionID % numShards), so messages
// from the same session are always processed in order by the same worker.
// Each shard has exactly one dedicated goroutine.
//
// When a shard queue is full the message is dropped and logged instead of
// blocking the caller's read loop.
type GoroutinePool struct {
	cfg          poolConfig
	shards       []chan IMessage
	stopCh       chan struct{}
	logger       ilog.ILogger
	eventHandler IEventHandler
	recycleMsg   func(IMessage)
}

func newGoroutinePool(cfg poolConfig, logger ilog.ILogger, handler IEventHandler, recycleMsg func(IMessage)) *GoroutinePool {
	p := &GoroutinePool{
		cfg:          cfg,
		stopCh:       make(chan struct{}),
		logger:       logger,
		eventHandler: handler,
		recycleMsg:   recycleMsg,
	}
	if cfg.queueSize == 0 {
		// queueSize=0: messages handled inline by the session's own goroutine.
		return p
	}
	numShards := int(cfg.maxIdleNum)
	if numShards == 0 {
		numShards = runtime.NumCPU()
	}
	p.shards = make([]chan IMessage, numShards)
	for i := range p.shards {
		p.shards[i] = make(chan IMessage, cfg.queueSize)
		go p.worker(p.shards[i])
	}
	return p
}

// push routes msg to the shard determined by session ID.
// If the shard queue is full the message is dropped and a warning is logged.
func (p *GoroutinePool) push(msg IMessage) {
	n := uint64(len(p.shards))
	if n == 0 {
		return
	}
	idx := msg.From().ID() % n
	select {
	case p.shards[idx] <- msg:
	default:
		p.logger.Errorf("gonet: shard[%d] queue full, dropping msg %d from session %d",
			idx, msg.ID(), msg.From().ID())
		p.recycleMsg(msg)
	}
}

func (p *GoroutinePool) worker(queue chan IMessage) {
	defer func() {
		if r := recover(); r != nil {
			p.logger.Errorf("gonet worker panic: %v\n%s", r, debug.Stack())
			// Restart on the same shard queue to maintain worker count.
			go p.worker(queue)
		}
	}()
	for {
		select {
		case <-p.stopCh:
			return
		case msg := <-queue:
			p.eventHandler.OnMessage(msg)
			p.recycleMsg(msg)
		}
	}
}

// Stop signals all shard workers to exit gracefully.
func (p *GoroutinePool) Stop() {
	close(p.stopCh)
}
