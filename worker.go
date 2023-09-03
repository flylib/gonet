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

type IEvent interface {
	Session() ISession
	Message() IMessage
}

type event struct {
	session ISession
	message IMessage
}

func (e event) Session() ISession {
	return e.session
}

func (e event) Message() IMessage {
	return e.message
}

// 处理池
type BeeWorkerPool struct {
	*Context
	//当前池协程数量(池大小)
	size int32
	//接受处理消息通道
	rcvMsgCh, events chan IEvent
	//创建协程通知
	createWorkerCh chan int
	//消息溢满通知
	overflowNotifyCh chan int
	//消息缓存
	msgCache IEventCache
}

// 初始化协程池
func createBeeWorkerPool(c *Context, size int32, msgCache IEventCache) (pool BeeWorkerPool) {
	pool = BeeWorkerPool{
		Context:          c,
		createWorkerCh:   make(chan int),
		overflowNotifyCh: make(chan int, 1),
		rcvMsgCh:         make(chan IEvent),
		events:           make(chan IEvent, receiveQueueSize),
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

func (self *BeeWorkerPool) createBeeWorker(count int32) {
	if count <= 0 {
		count = 1
	}
	for i := int32(0); i < count; i++ {
		self.createWorkerCh <- 1
	}
}

func (self *BeeWorkerPool) handle(e IEvent) {
	if len(self.events) >= receiveQueueSize {
		self.msgCache.Push(e)
		if len(self.overflowNotifyCh) < 1 {
			self.overflowNotifyCh <- 1
		}
	} else {
		self.events <- e
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
				for e := range self.events {
					if f, ok := self.Context.mMsgHooks[e.Message().ID()]; ok {
						f(e.Session(), e.Message())
					}
				}
			}()
		}
	}()
	//消息缓存处理
	go func() {
		for {
			if e := self.msgCache.Pop(); e != nil {
				self.events <- e
			} else {
				<-self.overflowNotifyCh
			}
		}
	}()
}
