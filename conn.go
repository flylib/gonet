package gonet

import (
	"reflect"
	"sync"
	"sync/atomic"
)

//////////////////////////////
////    Connection POOL   ////
//////////////////////////////

//链接
type Conn interface {
	//ID
	ID() uint64
	//原始套接字
	Socket() interface{}
	//断开
	Close()
	//发送消息
	Send(msg interface{}) error
	//设置键值对，存储关联数据
	Store(key string, value interface{})
	//获取键值对
	Load(key string) (value interface{}, ok bool)
}

//会话管理
type manager struct {
	incr  uint64    //流水号
	coons sync.Map  //所有链接
	pool  sync.Pool //conn临时对象池
}

func GetConn(id uint64) (Conn, bool) {
	value, ok := mgr.coons.Load(id)
	if ok {
		return value.(Conn), ok
	}
	return nil, false
}

func GetConnFromPool() Session {
	conn := mgr.pool.Get()
	globalLock.Lock()
	mgr.incr++
	conn.(interface{ setID(id uint64) }).setID(mgr.incr)
	globalLock.Unlock()
	mgr.Store(mgr.incr, conn)
	session := conn.(Session)
	//notify
	msg := &Msg{
		Session: session,
		SceneID: GetMsgSceneID(MsgIDSessionConnect),
		ID:      MsgIDSessionConnect,
		Data:    &msgSessionConnect,
	}
	PushWorkerPool(msg)
	return session
}

func RecycleConn(conn Conn) {
	conn.Close()
	mgr.coons.Delete(conn.ID())
	mgr.pool.Put(conn)
}

func SessionCount() int {
	sum := 0
	mgr.Range(func(key, value interface{}) bool {
		sum++
		return true
	})
	return sum
}

//广播
func Broadcast(msg interface{}) {
	mgr.Range(func(_, item interface{}) bool {
		item.(Session).Send(msg)
		return true
	})
}
