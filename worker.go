package gonet

import (
	"math"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"time"
)

type poolConfig struct {
	queueSize  int32
	maxNum     int32
	maxIdleNum int32
}

// Lightweight goroutine pool
type GoroutinePool struct {
	cfg poolConfig
	*Context
	curWorkingNum     int32
	cacheQueueSize    int
	queue             chan IMessage
	addRoutineChannel chan bool
}

func newGoroutinePool(ctx *Context) *GoroutinePool {
	if ctx.poolCfg.maxIdleNum == 0 {
		ctx.poolCfg.maxIdleNum = int32(runtime.NumCPU())
	}
	if ctx.poolCfg.queueSize == 0 {
		ctx.poolCfg.queueSize = 64
	}

	pool := &GoroutinePool{
		Context:           ctx,
		cfg:               ctx.poolCfg,
		addRoutineChannel: make(chan bool),
		queue:             make(chan IMessage, ctx.poolCfg.queueSize),
	}

	go pool.run()
	pool.ascRoutine(pool.cfg.maxIdleNum)
	return pool
}

func (b *GoroutinePool) ascRoutine(count int32) {
	if count <= 0 {
		count = 1
	}
	for i := int32(0); i < count; i++ {
		b.addRoutineChannel <- true
	}
}

func (b *GoroutinePool) descRoutine(count int32) {
	if count <= 0 {
		count = 1
	}
	for i := int32(0); i < count; i++ {
		//b.queue <- newErrorMessage(nil)
	}
}

func (b *GoroutinePool) run() {
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
				b.Context.messageHandler(e)
			}
		}()
	}
}

func (b *GoroutinePool) monitor() {
	if b.cfg.maxNum == 0 {
		return
	}
	tick := time.Tick(time.Second * 30)
	var preCount int
	for range tick {
		curCount := len(b.queue)
		between := curCount - preCount
		if between > 0 {
			count := math.Abs(float64(between) / float64(b.cacheQueueSize) * float64(b.cfg.maxNum))
			b.ascRoutine(int32(count))
		} else if preCount > 0 && curCount == 0 {
			curWorkingNum := atomic.LoadInt32(&b.curWorkingNum)
			if curWorkingNum > b.cfg.maxIdleNum {
				b.descRoutine(1)
			}
		}
	}
}
