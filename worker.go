package goNet

import (
	"github.com/astaxie/beego/logs"
	"runtime"
	"sync/atomic"
)

//////////////////////////
////    WORKER POOL   ////
//////////////////////////

var workers WorkerPool

//处理池
type WorkerPool struct {
	//事件管道
	eventChannel chan Event
	//因提交事件阻塞的协程数量
	blockingNum int32
	//当前池协程数量(池大小)
	size int32
	//创建协程通知
	createNotify chan interface{}
	//异常处理函数
	panicHandler func(interface{})
	//池限制大小
	//默认 runtime.NumCPU() * 10
	maxPoolSize int32
}

//初始化协程池
func InitWorkerPool(panicHandler func(interface{}), eventChannelSize int) {
	workers = WorkerPool{
		createNotify: make(chan interface{}),
		panicHandler: panicHandler,
		maxPoolSize:  int32(runtime.NumCPU() * 10),
	}
	if eventChannelSize < 1 {
		workers.eventChannel = make(chan Event) //无缓存通道
	} else {
		workers.eventChannel = make(chan Event, eventChannelSize) //有缓存通道
	}
	workers.run()
	workers.createWorker(1)
}

func (w *WorkerPool) incBlocking() {
	atomic.AddInt32(&w.blockingNum, 1)
}
func (w *WorkerPool) decBlocking() {
	atomic.AddInt32(&w.blockingNum, -1)
}
func (w *WorkerPool) incPoolSize() {
	atomic.AddInt32(&w.size, 1)
}
func (w *WorkerPool) decPoolSize() {
	atomic.AddInt32(&w.size, -1)
}
func (w *WorkerPool) createWorker(count int) {
	go func() {
		for i := 0; i < count; i++ {
			w.createNotify <- Event{eventType: EventWorkerAdd}
		}
	}()
}

func (w *WorkerPool) destroyWorker() {
	w.createNotify <- Event{eventType: EventWorkerExit}
}

//处理事件
func (w *WorkerPool) handling(e Event) {
	w.incBlocking()
	//todo 按需调整池大小,不做精确控制
	if w.blockingNum > int32(runtime.NumCPU()) && w.size < w.maxPoolSize {
		w.createWorker(1)
	} else if w.size > int32(runtime.NumCPU()) {
		w.destroyWorker()
	}
	w.eventChannel <- e
	w.decBlocking()
}

//运行
func (w *WorkerPool) run() {
	go func() {
		for range w.createNotify {
			w.incPoolSize()
			go func() {
				//panic handling
				defer func() {
					w.decPoolSize()
					if info := recover(); info != nil {
						if w.panicHandler != nil {
							w.panicHandler(info)
						} else {
							logs.Error("worker exits from a panic: %v\n", info)
							var buf [4096]byte
							n := runtime.Stack(buf[:], false)
							logs.Error("worker exits from panic: %s\n", string(buf[:n]))
						}
					}
				}()
				logs.Info("new worker wait")
				for e := range w.eventChannel {
					logs.Info("new msg")
					if e.eventType == EventWorkerExit {
						w.decPoolSize()
						return
					}
					e.route.OnMsg(e.from, e.data)
				}
			}()
		}
	}()
}

//处理事件
func HandleEvent(event Event) {
	workers.handling(event)
}
