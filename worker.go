package goNet

import (
	"github.com/astaxie/beego/logs"
	"runtime"
	"sync/atomic"
)

//////////////////////////
////     WORKER       ////
//////////////////////////

var wPool workerPool

//事件分类
const (
	EventNetWorkIO  int8 = iota //网络io
	EventWorkerExit             //worker退出
	EventWorkerAdd              //新增worker
)

//事件
type Event struct {
	Kind   int8        //事件分类
	From   Session     //来自
	Router Route       //路由(处理器)
	Msg    interface{} //消息-事件内容
}

//处理池
type workerPool struct {
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
	//默认100
	maxPoolSize int32
}

//初始化协程池
func InitWorkerPool(panicHandler func(interface{}), eventChannelSize int) {
	wPool = workerPool{
		//eventChannel:  make(chan Event),
		createNotify: make(chan interface{}),
		panicHandler: panicHandler,
		maxPoolSize:  int32(runtime.NumCPU() * 10),
	}
	if eventChannelSize < 1 {
		wPool.eventChannel = make(chan Event)
	} else {
		wPool.eventChannel = make(chan Event, eventChannelSize)
	}
	wPool.run()
	wPool.createWorker(1)
}

func (w *workerPool) incBlocking() {
	atomic.AddInt32(&w.blockingNum, 1)
	logs.Info("p.blockingNum=", w.blockingNum)
}
func (w *workerPool) decBlocking() {
	atomic.AddInt32(&w.blockingNum, -1)
}

func (w *workerPool) incPoolSize() {
	atomic.AddInt32(&w.size, 1)
	logs.Info("p.size=", w.size)
}
func (w *workerPool) decPoolSize() {
	atomic.AddInt32(&w.size, -1)
}
func (w *workerPool) createWorker(count int) {
	go func() {
		for i := 0; i < count; i++ {
			w.createNotify <- EventWorkerAdd
		}
	}()
}

func (w *workerPool) destroyWorker() {
	w.createNotify <- Event{Kind: EventWorkerExit}
}

//处理事件
func (w *workerPool) handling(e Event) {
	w.incBlocking()
	//todo 按需调整池大小
	if w.blockingNum > int32(runtime.NumCPU()) {
		//不做精确控制
		if w.size < w.maxPoolSize {
			w.createWorker(1)
		}
	} else {
		//不做精确控制
		if w.size > int32(runtime.NumCPU()) {
			w.destroyWorker()
		}
	}
	w.eventChannel <- e
	w.decBlocking()
}

//运行
func (w *workerPool) run() {
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
					if e.Kind == EventWorkerExit {
						w.decPoolSize()
						return
					}
					e.Router.OnMsg(e.From, e.Msg)
				}
			}()
		}
	}()
}

//提交到协程池
func CommitWorkerPool(event Event) {
	wPool.handling(event)
}
