package gonet

import (
	"log"
	"runtime"
	"sync/atomic"
)

//////////////////////////
////    WORKER POOL   ////
//////////////////////////

const (
	receiveQueueSize = 512 //默认接收队列大小
)

//处理池
type WorkerPool struct {
	//当前池协程数量(池大小)
	size int32
	//接受处理消息通道
	sessionCh, handleSessionCh chan *Session
	//创建协程通知
	createWorkerCh chan int
	//消息溢满通知
	overflowNotifyCh chan int
	//消息缓存
	sessionCache SessionCache
}

//初始化协程池
func createWorkerPool(size int32, sessionCache SessionCache) (pool WorkerPool) {
	pool = WorkerPool{
		createWorkerCh:   make(chan int),
		overflowNotifyCh: make(chan int, 1),
		sessionCh:        make(chan *Session),
		handleSessionCh:  make(chan *Session, receiveQueueSize),
		sessionCache:     sessionCache,
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
func (w *WorkerPool) handle(msg *Session) {
	if len(w.handleSessionCh) >= receiveQueueSize {
		w.sessionCache.Push(msg)
		if len(w.overflowNotifyCh) < 1 {
			w.overflowNotifyCh <- 1
		}
	} else {
		w.handleSessionCh <- msg
	}
}

//运行
func (w *WorkerPool) run() {
	go func() {
		for msg := range w.sessionCh {
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
						var buf [4096]byte
						n := runtime.Stack(buf[:], false)
						log.Printf("worker exits from panic: %s\n", string(buf[:n]))
					}
					w.createWorkerCh <- 1
				}()
				for s := range w.handleSessionCh {
					if f, ok := sys.mHandlers[s.Msg.ID]; ok {
						f(s)
					}
				}
			}()
		}
	}()
	//消息缓存处理
	go func() {
		for {
			if e := w.sessionCache.Pop(); e != nil {
				w.handleSessionCh <- e
			} else {
				<-w.overflowNotifyCh
			}
		}
	}()
}
