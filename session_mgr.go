package gonet

import (
	"sync"
)

//会话管理
type SessionManager struct {
	sync.RWMutex
	incr     uint64    //流水号
	sessions sync.Map  //所有链接
	pool     sync.Pool //临时对象池
}

func (s *SessionManager) store(id uint64, session interface{}) {
	session.(interface{ setID(id uint64) }).setID(id)
	sys.sessions.Store(id, session)
}

func (s *SessionManager) del(id uint64) {
	s.sessions.Delete(id)
}
