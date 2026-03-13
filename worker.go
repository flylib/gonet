package gonet

import (
	"runtime"
	"runtime/debug"
	"sync/atomic"

	ilog "github.com/flylib/interface/log"
)

type poolConfig struct {
	queueSize  int32
	maxNum     int32
	maxIdleNum int32
}

// GoroutinePool is a lightweight goroutine pool for processing messages.
type GoroutinePool struct {
	cfg          poolConfig
	queue        chan IMessage
	addCh        chan struct{}
	stopCh       chan struct{}
	curWorkers   int32
	logger       ilog.ILogger
	eventHandler IEventHandler
}

func newGoroutinePool(cfg poolConfig, logger ilog.ILogger, handler IEventHandler) *GoroutinePool {
	p := &GoroutinePool{
		cfg:          cfg,
		stopCh:       make(chan struct{}),
		logger:       logger,
		eventHandler: handler,
	}
	if cfg.queueSize == 0 {
		// No pool: messages are processed inline by the session's own goroutine.
		return p
	}
	if cfg.maxIdleNum == 0 {
		cfg.maxIdleNum = int32(runtime.NumCPU())
	}
	p.queue = make(chan IMessage, cfg.queueSize)
	p.addCh = make(chan struct{}, cfg.maxIdleNum+1)
	go p.supervisor()
	p.addWorkers(cfg.maxIdleNum)
	return p
}

func (p *GoroutinePool) push(msg IMessage) {
	p.queue <- msg
}

// addWorkers signals the supervisor to start n new workers.
func (p *GoroutinePool) addWorkers(n int32) {
	for i := int32(0); i < n; i++ {
		select {
		case p.addCh <- struct{}{}:
		default:
		}
	}
}

// supervisor listens for add-worker signals and enforces the maxNum cap.
func (p *GoroutinePool) supervisor() {
	for range p.addCh {
		if p.cfg.maxNum > 0 && atomic.LoadInt32(&p.curWorkers) >= p.cfg.maxNum {
			continue
		}
		atomic.AddInt32(&p.curWorkers, 1)
		go p.worker()
	}
}

func (p *GoroutinePool) worker() {
	defer func() {
		atomic.AddInt32(&p.curWorkers, -1)
		if r := recover(); r != nil {
			p.logger.Errorf("gonet worker panic: %v\n%s", r, debug.Stack())
			// restart to maintain pool size after panic
			p.addWorkers(1)
		}
	}()

	for {
		select {
		case <-p.stopCh:
			return
		case msg := <-p.queue:
			p.eventHandler.OnMessage(msg)
			if m, ok := msg.(*message); ok {
				recycleMessage(m)
			}
		}
	}
}

// Stop signals all workers to exit gracefully.
func (p *GoroutinePool) Stop() {
	close(p.stopCh)
}
