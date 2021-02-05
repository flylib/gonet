package goNet

import (
	"github.com/astaxie/beego/logs"
	"runtime"
	"sync/atomic"
	"time"
)

//////////////////////////
////    WORKER POOL   ////
//////////////////////////

var workers WorkerPool

//处理池
type WorkerPool struct {
	//事件管道
	reciveCh chan *Msg
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
func NewWorkerPool(panicHandler func(interface{})) (pool WorkerPool) {
	pool = WorkerPool{
		createNotify: make(chan interface{}),
		panicHandler: panicHandler,
		size:         int32(runtime.NumCPU()),
		reciveCh:     make(chan *Msg),
	}
	pool.run()
	pool.createWorker(1)
	pool.tick()
	return pool
}

func (w *WorkerPool) tick() {
	go func() {
		minDecPoolSize := int32(float64(w.maxPoolSize) * 0.5)
		count := 0
		for {
			count++
			//每隔1分检查一下
			time.Sleep(time.Minute)
			//两倍扩容速度
			if w.blockingNum > 10 {
				curSize := w.size
				if curSize*2 < w.maxPoolSize {
					w.createWorker(curSize)
				} else {
					w.createWorker(atomic.LoadInt32(&w.maxPoolSize) - atomic.LoadInt32(&w.size))
				}
			}
			//每30分钟检查一次缩容，0.85倍缩容
			if count > 30 {
				count = 0
				curPoolSize := atomic.LoadInt32(&w.size)
				if curPoolSize > minDecPoolSize {
					destroyPoolSize := int32(float64(curPoolSize) * 0.85)
					w.destroyWorker(destroyPoolSize)
				}
			}
		}
	}()
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

func (w *WorkerPool) createWorker(count int32) {
	for i := int32(0); i < count; i++ {
		w.createNotify <- i
	}
}
func (w *WorkerPool) destroyWorker(count int32) {
	for i := int32(0); i < count; i++ {
		w.reciveCh <- &Msg{
			ID: MsgIDDecPoolSize,
		}
	}
}

//运行
func (w *WorkerPool) run() {
	go func() {
		for range w.createNotify {
			if w.size >= w.maxPoolSize {
				continue
			}
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
				for msg := range w.reciveCh {
					if msg.ID == MsgIDDecPoolSize {
						w.decPoolSize()
						return
					}
					scene := msg.GetScene(msg.SceneID)
					if scene != nil {
						scene.Handler(msg)
					}
				}
			}()
		}
	}()
}

//放进工作池
func PushWorkerPool(msg *Msg) {
	workers.incBlocking()
	workers.reciveCh <- msg
	workers.decBlocking()
}
