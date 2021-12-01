package gonet

import (
	"github.com/astaxie/beego/logs"
	"runtime"
	"sync/atomic"
	"time"
)

//////////////////////////
////    WORKER POOL   ////
//////////////////////////

const (
	receiveQueueSize = 128 //接收队列大小
)

var workers WorkerPool

func initWorkerPool(option Option) {
	workers = NewWorkerPool(option.maxWorkerPoolSize)
}

//处理池
type WorkerPool struct {
	//当前池协程数量(池大小)
	//池限制大小
	size, maxSize int32
	//接受消息通道
	receiveMsgCh chan *Message
	//创建协程通知
	createWorkerCh chan int
}

//初始化协程池
func NewWorkerPool(maxPoolSize int32) (pool WorkerPool) {
	if maxPoolSize < 1 {
		maxPoolSize = 1
	}
	pool = WorkerPool{
		createWorkerCh: make(chan int),
		maxSize:        maxPoolSize,
		receiveMsgCh:   make(chan *Message, receiveQueueSize),
	}
	pool.run()
	pool.createWorker(1)
	pool.tick()
	return pool
}

func (w *WorkerPool) tick() {
	go func() {
		minDecPoolSize := int32(float64(w.maxSize) * 0.5)
		count := 0
		for {
			count++
			//每隔1分检查一下
			time.Sleep(time.Minute)
			//两倍扩容速度
			if len(w.receiveMsgCh) == cap(w.receiveMsgCh) && atomic.LoadInt32(&w.size) < atomic.LoadInt32(&w.maxSize) {
				curSize := w.size
				if curSize*2 < w.maxSize {
					w.createWorker(curSize)
				} else {
					w.createWorker(w.maxSize - w.size)
				}
			}
			//每30分钟检查一次缩容，0.85倍缩容
			if count > 30 {
				count = 0
				if len(w.receiveMsgCh) < 1 {
					curPoolSize := atomic.LoadInt32(&w.size)
					if curPoolSize > minDecPoolSize {
						destroyPoolSize := int32(float64(curPoolSize) * 0.85)
						w.destroyWorker(destroyPoolSize)
					}
				}
			}
		}
	}()
}

func (w *WorkerPool) incPoolSize() {
	atomic.AddInt32(&w.size, 1)
}
func (w *WorkerPool) decPoolSize() {
	atomic.AddInt32(&w.size, -1)
}

func (w *WorkerPool) createWorker(count int32) {
	for i := int32(0); i < count; i++ {
		w.createWorkerCh <- 1
	}
}
func (w *WorkerPool) destroyWorker(count int32) {
	for i := int32(0); i < count; i++ {
		w.receiveMsgCh <- &Message{
			ID: MsgIDDecPoolSize,
		}
	}
}

//运行
func (w *WorkerPool) run() {
	go func() {
		for range w.createWorkerCh {
			if w.size >= w.maxSize {
				continue
			}
			w.incPoolSize()
			go func() {
				//panic handling
				defer func() {
					w.decPoolSize()
					if info := recover(); info != nil {
						var buf [4096]byte
						n := runtime.Stack(buf[:], false)
						logs.Error("worker exits from panic: %s\n", string(buf[:n]))
					}
				}()
				for msg := range w.receiveMsgCh {
					if msg.ID == MsgIDDecPoolSize {
						w.decPoolSize()
						return
					}
					scene := msg.GetScene(msg.SceneID)
					if scene != nil {
						scene.Handler(msg)
						continue
					}
					scene = getCommonScene(msg.SceneID)
					if scene != nil {
						scene.Handler(msg)
					}
				}
			}()
		}
	}()
}

//放进工作池
func PushWorkerPool(msg *Message) {
	workers.receiveMsgCh <- msg
}
