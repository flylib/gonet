package gonet

import (
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	sys           System       //系统
	transportType reflect.Type //传输协议类型
)

type System struct {
	SessionManager
	//msgID:msgType
	msgTypes map[MessageID]reflect.Type
	//msgType:msgID
	msgIDs map[reflect.Type]MessageID
}

func init() {
	sys = System{
		SessionManager: SessionManager{
			pool: sync.Pool{
				New: func() interface{} {
					return reflect.New(transportType).Interface()
				},
			},
		},
	}
}

//获取会话
func GetSession(id uint64) (Session, bool) {
	value, ok := sys.sessions.Load(id)
	if ok {
		return value.(Session), ok
	}
	return nil, false
}

//创建会话
func CreateSession() Session {
	obj := sys.pool.Get()
	sys.store(atomic.AddUint64(&sys.incr, 1), obj)
	session := obj.(Session)
	return session
}

//回收会话对象
func RecycleSession(session Session) {
	//关闭
	session.Close()
	//删除
	sys.del(session.ID())
	//回收
	sys.pool.Put(session)
}

//统计会话数量
func SessionCount() int {
	sum := 0
	sys.sessions.Range(func(key, value interface{}) bool {
		sum++
		return true
	})
	return sum
}

//广播会话
func Broadcast(msg interface{}) {
	sys.sessions.Range(func(_, item interface{}) bool {
		session, ok := item.(Session)
		if ok {
			session.Send(msg)
		}
		return true
	})
}
