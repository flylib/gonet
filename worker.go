package gonet

import (
	"math"
	"runtime/debug"
	"sync/atomic"
	"time"
)

/*----------------------------------------------------------------
					////////////////////////////
					////  BEE WORKER POOL   ////
					////////////////////////////
----------------------------------------------------------------*/

const (
	receiveQueueSize = 512 //默认接收队列大小
)

type BeeWorkerPool struct {
	*AppContext
	curWorkingNum   int32
	maxWorkerNum    int32
	idleWorkerNum   int32
	cacheQueueSize  int
	queue           chan IMessage
	addWorkerNotify chan bool
}

func maxBeeWorkers(num int32) func(pool *BeeWorkerPool) {
	return func(pool *BeeWorkerPool) {
		pool.maxWorkerNum = num
	}
}

func allowIdleBeesWorkers(num int32) func(pool *BeeWorkerPool) {
	return func(pool *BeeWorkerPool) {
		pool.idleWorkerNum = num
	}
}

func setQueueSize(num int) func(pool *BeeWorkerPool) {
	if num < 0 {
		num = 0
	}
	return func(pool *BeeWorkerPool) {
		pool.cacheQueueSize = num
		pool.queue = make(chan IMessage, num)
	}
}

func newBeeWorkerPool(c *AppContext, options ...func(pool *BeeWorkerPool)) (pool BeeWorkerPool) {
	pool = BeeWorkerPool{
		AppContext:      c,
		addWorkerNotify: make(chan bool),
		queue:           make(chan IMessage, receiveQueueSize),
	}

	for _, option := range options {
		option(&pool)
	}

	go pool.run()
	pool.addBeeWorker(pool.idleWorkerNum)
	return pool
}

func (b *BeeWorkerPool) addBeeWorker(count int32) {
	if count <= 0 {
		count = 1
	}
	for i := int32(0); i < count; i++ {
		b.addWorkerNotify <- true
	}
}

func (b *BeeWorkerPool) descBeeWorker(count int32) {
	if count <= 0 {
		count = 1
	}
	for i := int32(0); i < count; i++ {
		b.queue <- newInvalidMessage()
	}
}

func (b *BeeWorkerPool) run() {
	for range b.addWorkerNotify {
		if b.curWorkingNum >= b.maxWorkerNum {
			continue
		}
		atomic.AddInt32(&b.curWorkingNum, 1)
		go func() {
			// panic handling
			defer func() {
				atomic.AddInt32(&b.curWorkingNum, -1)
				if err := recover(); err != nil {
					b.Errorf("panic error:%s\n%s", err, debug.Stack())
					b.addBeeWorker(1)
				}
			}()

			// message handling
			for msg := range b.queue {
				if f, ok := b.AppContext.GetMessageHandler(msg.ID()); ok {
					f(msg)
				} else {
					break //release go routine
				}
			}
		}()
	}
}

func (b *BeeWorkerPool) monitor() {
	tick := time.Tick(time.Second * 30)
	var preCount int
	for range tick {
		curCount := len(b.queue)
		between := curCount - preCount
		if between > 0 {
			count := math.Abs(float64(between) / float64(b.cacheQueueSize) * float64(b.maxWorkerNum))
			b.addBeeWorker(int32(count))
		} else if preCount > 0 && curCount == 0 {
			curWorkingNum := atomic.LoadInt32(&b.curWorkingNum)
			if curWorkingNum > b.idleWorkerNum {
				b.descBeeWorker(1)
			}
		}
	}
}
