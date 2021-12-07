package gonet

import (
	"container/list"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
)

//////////////////////////
////    WORKER POOL   ////
//////////////////////////

const (
	receiveQueueSize = 1024 //默认接收队列大小
)

//处理池
type WorkerPool struct {
	sync.Mutex
	//当前池协程数量(池大小)
	size int32
	//接受消息通道
	receiveMsgCh chan *Message
	//创建协程通知
	createWorkerCh chan int
	//消息溢满通知
	overflowNotifyCh chan int
	//消息缓存
	msgList *list.List
}

//初始化协程池
func newWorkerPool(size int32) (pool WorkerPool) {
	pool = WorkerPool{
		createWorkerCh:   make(chan int),
		overflowNotifyCh: make(chan int, 1),
		receiveMsgCh:     make(chan *Message, receiveQueueSize),
		msgList:          list.New(),
	}
	pool.run()
	pool.createWorker(size)
	return pool
}

func (w *WorkerPool) incPoolSize() {
	atomic.AddInt32(&w.size, 1)
}
func (w *WorkerPool) decPoolSize() {
	atomic.AddInt32(&w.size, -1)
}

func (w *WorkerPool) createWorker(count int32) {
	if count == 0 {
		count = 1
	}
	for i := int32(0); i < count; i++ {
		w.createWorkerCh <- 1
	}
}
func (w *WorkerPool) handle(msg *Message) {
	if len(sys.workers.receiveMsgCh) >= receiveQueueSize {
		w.Lock()
		defer w.Unlock()
		w.msgList.PushFront(msg)
		if len(w.overflowNotifyCh) < 1 {
			w.overflowNotifyCh <- 1
		}
	} else {
		sys.workers.receiveMsgCh <- msg
	}
}

//运行
func (w *WorkerPool) run() {
	go func() {
		for range w.createWorkerCh {
			w.incPoolSize()
			go func() {
				//panic handling
				defer func() {
					w.decPoolSize()
					if err := recover(); err != nil {
						var buf [4096]byte
						n := runtime.Stack(buf[:], false)
						log.Printf("worker exits from panic: %s\n", string(buf[:n]))
					}
					w.createWorkerCh <- 1
				}()
				for msg := range w.receiveMsgCh {
					if f, ok := sys.mHandlers[msg.ID]; ok {
						f(msg)
					}
				}
			}()
		}
	}()
	//消息缓存处理
	go func() {
		for {
			if e := w.msgList.Back(); e != nil {
				msg := e.Value.(*Message)
				w.receiveMsgCh <- msg
				w.msgList.Remove(e)
			} else {
				<-w.overflowNotifyCh
			}
		}
	}()
}
