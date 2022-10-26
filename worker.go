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

//处理池
type BeeWorkerPool struct {
	//当前池协程数量(池大小)
	size int32
	//接受处理消息通道
	rcvMsgCh, handleMsgCh chan *Message
	//创建协程通知
	createWorkerCh chan int
	//消息溢满通知
	overflowNotifyCh chan int
	//消息缓存
	msgCache MessageCache
}

//初始化协程池
func createBeeWorkerPool(size int32, msgCache MessageCache) (pool BeeWorkerPool) {
	pool = BeeWorkerPool{
		createWorkerCh:   make(chan int),
		overflowNotifyCh: make(chan int, 1),
		rcvMsgCh:         make(chan *Message),
		handleMsgCh:      make(chan *Message, receiveQueueSize),
		msgCache:         msgCache,
	}
	pool.run()
	pool.createBeeWorker(size)
	return pool
}

func (w *BeeWorkerPool) incPoolSize() {
	atomic.AddInt32(&w.size, 1)
}
func (w *BeeWorkerPool) decPoolSize() {
	atomic.AddInt32(&w.size, -1)
}

func (w *BeeWorkerPool) createBeeWorker(count int32) {
	if count <= 0 {
		count = 1
	}
	for i := int32(0); i < count; i++ {
		w.createWorkerCh <- 1
	}
}

func (w *BeeWorkerPool) handle(msg *Message) {
	if len(w.handleMsgCh) >= receiveQueueSize {
		w.msgCache.Push(msg)
		if len(w.overflowNotifyCh) < 1 {
			w.overflowNotifyCh <- 1
		}
	} else {
		w.handleMsgCh <- msg
	}
}

//运行
func (w *BeeWorkerPool) run() {
	go func() {
		for msg := range w.rcvMsgCh {
			w.handle(msg)
		}
	}()
	go func() {
		for range w.createWorkerCh {
			w.incPoolSize()
			go func() {
				//panic handling
				defer func() {
					w.decPoolSize()
					if err := recover(); err != nil {
						fmt.Fprintf(os.Stderr, "panic error:%s\n%s", err, debug.Stack())
					}
					w.createWorkerCh <- 1
				}()
				for msg := range w.handleMsgCh {
					if f, ok := ctx.mMsgHooks[msg.ID]; ok {
						f(msg)
					}
				}
			}()
		}
	}()
	//消息缓存处理
	go func() {
		for {
			if e := w.msgCache.Pop(); e != nil {
				w.handleMsgCh <- e
			} else {
				<-w.overflowNotifyCh
			}
		}
	}()
}
