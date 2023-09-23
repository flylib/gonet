package gonet

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync/atomic"
)

/*----------------------------------------------------------------
					////////////////////////////
					////  BEE WORKER POOL   ////
					////////////////////////////
----------------------------------------------------------------*/

const (
	receiveQueueSize = 512 //默认接收队列大小
)

// 处理池
type BeeWorkerPool struct {
	*AppContext
	//当前池协程数量(池大小)
	size int32
	//接受处理消息通道
	rcvMsgCh, handingCh chan IMessage
	//创建协程通知
	createWorkerCh chan int
	//消息溢满通知
	overflowNotifyCh chan int
	//消息缓存
	msgCache IMessageCache
}

// 初始化协程池
func createBeeWorkerPool(c *AppContext, size int, msgCache IMessageCache) (pool BeeWorkerPool) {
	pool = BeeWorkerPool{
		AppContext:       c,
		createWorkerCh:   make(chan int),
		overflowNotifyCh: make(chan int, 1),
		rcvMsgCh:         make(chan IMessage),
		handingCh:        make(chan IMessage, receiveQueueSize),
		msgCache:         msgCache,
	}
	pool.run()
	pool.createBeeWorker(size)
	return pool
}

func (self *BeeWorkerPool) incPoolSize() {
	atomic.AddInt32(&self.size, 1)
}
func (self *BeeWorkerPool) decPoolSize() {
	atomic.AddInt32(&self.size, -1)
}

func (self *BeeWorkerPool) createBeeWorker(count int) {
	if count <= 0 {
		count = 1
	}
	for i := 0; i < count; i++ {
		self.createWorkerCh <- 1
	}
}

func (self *BeeWorkerPool) handle(e IMessage) {
	if len(self.handingCh) >= receiveQueueSize {
		self.msgCache.Push(e)
		if len(self.overflowNotifyCh) < 1 {
			self.overflowNotifyCh <- 1
		}
	} else {
		self.handingCh <- e
	}
}

// 运行
func (self *BeeWorkerPool) run() {
	go func() {
		for msg := range self.rcvMsgCh {
			self.handle(msg)
		}
	}()
	go func() {
		for range self.createWorkerCh {
			self.incPoolSize()
			go func() {
				//panic handling
				defer func() {
					self.decPoolSize()
					if err := recover(); err != nil {
						fmt.Fprintf(os.Stderr, "panic error:%s\n%s", err, debug.Stack())
					}
					self.createWorkerCh <- 1
				}()
				for msg := range self.handingCh {
					if f, ok := self.AppContext.mMsgHooks[msg.ID()]; ok {
						f(msg)
					}
				}
			}()
		}
	}()
	//消息缓存处理
	go func() {
		for {
			if e := self.msgCache.Pop(); e != nil {
				self.handingCh <- e
			} else {
				<-self.overflowNotifyCh
			}
		}
	}()
}
