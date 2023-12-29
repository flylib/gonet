package gonet

import (
	"runtime"
	"runtime/debug"
	"sync/atomic"
)

type RuntimeConfig struct {
	queueSize  int32
	maxNum     int32
	maxIdleNum int32
}

// Lightweight goroutine pool
type AsyncRuntime struct {
	*Context
	cfg               RuntimeConfig
	curWorkingNum     int32
	cacheQueueSize    int
	queue             chan Message
	addRoutineChannel chan bool
}

func newAsyncRuntime(ctx *Context) *AsyncRuntime {
	if ctx.poolCfg.maxIdleNum == 0 {
		ctx.poolCfg.maxIdleNum = int32(runtime.NumCPU())
	}
	if ctx.poolCfg.queueSize == 0 {
		ctx.poolCfg.queueSize = 64
	}

	pool := &AsyncRuntime{
		Context:           ctx,
		cfg:               ctx.poolCfg,
		addRoutineChannel: make(chan bool),
		queue:             make(chan Message, ctx.poolCfg.queueSize),
	}

	go pool.run()
	pool.ascRoutine(pool.cfg.maxIdleNum)
	return pool
}

func (b *AsyncRuntime) ascRoutine(count int32) {
	if count <= 0 {
		count = 1
	}
	for i := int32(0); i < count; i++ {
		b.addRoutineChannel <- true
	}
}

func (b *AsyncRuntime) descRoutine(count int32) {
	if count <= 0 {
		count = 1
	}
	for i := int32(0); i < count; i++ {
		//b.queue <- newErrorMessage(nil)
	}
}

func (b *AsyncRuntime) run() {
	for range b.addRoutineChannel {
		if b.cfg.maxNum != 0 &&
			b.curWorkingNum >= b.cfg.maxNum {
			continue
		}
		atomic.AddInt32(&b.curWorkingNum, 1)
		go func() {
			// panic handling
			defer func() {
				atomic.AddInt32(&b.curWorkingNum, -1)
				if err := recover(); err != nil {
					b.ILogger.Errorf("panic error:%s\n%s", err, debug.Stack())
					b.ascRoutine(1)
				}
			}()

			// message handling
			for e := range b.queue {
				b.eventHandler.OnMessage(e)
			}
		}()
	}
}
