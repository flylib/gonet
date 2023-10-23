package gonet

import (
	"math"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"time"
)

type poolConfig struct {
	queueSize  uint32
	maxNum     uint32
	maxIdleNum uint32
}

// Lightweight goroutine pool
type GoroutinePool struct {
	*AppContext
	curWorkingNum, maxWorkingNum, maxIdleNum int32
	cacheQueueSize                           int
	queue                                    chan IMessage
	addRoutineChannel                        chan bool
}

func maxWorkingGoroutines(num int32) goroutinePoolOption {
	return func(pool *GoroutinePool) {
		pool.maxWorkingNum = num
	}
}

func maxIdleGoroutines(num int32) goroutinePoolOption {
	return func(pool *GoroutinePool) {
		pool.maxIdleNum = num
	}
}

func setQueueSize(num int) goroutinePoolOption {
	if num < 0 {
		num = 0
	}
	return func(pool *GoroutinePool) {
		pool.cacheQueueSize = num
		pool.queue = make(chan IMessage, num)
	}
}

func newGoroutinePool(c *AppContext, options ...goroutinePoolOption) *GoroutinePool {
	pool := &GoroutinePool{
		AppContext:        c,
		addRoutineChannel: make(chan bool),
		queue:             make(chan IMessage, defaultReceiveQueueSize),
		maxIdleNum:        int32(runtime.NumCPU()),
	}

	for _, opt := range options {
		option(pool)
	}

	go pool.run()
	pool.ascRoutine(pool.maxIdleNum)
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
		if b.maxWorkingNum != 0 &&
			b.curWorkingNum >= b.maxWorkingNum {
			continue
		}
		atomic.AddInt32(&b.curWorkingNum, 1)
		go func() {
			// panic handling
			defer func() {
				atomic.AddInt32(&b.curWorkingNum, -1)
				if err := recover(); err != nil {
					b.Errorf("panic error:%s\n%s", err, debug.Stack())
					b.ascRoutine(1)
				}
			}()

			// message handling
			for e := range b.queue {
				b.AppContext.callback(e)
			}
		}()
	}
}

func (b *GoroutinePool) monitor() {
	if b.maxWorkingNum == 0 {
		return
	}
	tick := time.Tick(time.Second * 30)
	var preCount int
	for range tick {
		curCount := len(b.queue)
		between := curCount - preCount
		if between > 0 {
			count := math.Abs(float64(between) / float64(b.cacheQueueSize) * float64(b.maxWorkingNum))
			b.ascRoutine(int32(count))
		} else if preCount > 0 && curCount == 0 {
			curWorkingNum := atomic.LoadInt32(&b.curWorkingNum)
			if curWorkingNum > b.maxIdleNum {
				b.descRoutine(1)
			}
		}
	}
}
