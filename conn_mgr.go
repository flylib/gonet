package gonet

import (
	"sync"
)

//会话管理
type ConnManager struct {
	sync.RWMutex
	incr        uint64    //流水号
	connections sync.Map  //所有链接
	pool        sync.Pool //临时对象池
}

func (s *ConnManager) store(id uint64, conn interface{}) {
	conn.(interface{ setID(id uint64) }).setID(id)
	s.connections.Store(id, conn)
}

func (s *ConnManager) del(id uint64) {
	s.connections.Delete(id)
}
