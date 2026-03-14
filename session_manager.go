package gonet

import (
	"sync"
	"sync/atomic"
)

// sessionManager manages alive sessions and the idle session pool.
type sessionManager[S SessionConstraint] struct {
	aliveNum int32
	serial   uint64
	alive    sync.Map
	pool     sync.Pool
}

func newSessionManager[S SessionConstraint](factory func() S) *sessionManager[S] {
	m := &sessionManager[S]{}
	m.pool.New = func() any { return factory() }
	return m
}

func (m *sessionManager[S]) addAlive(s S) {
	atomic.AddInt32(&m.aliveNum, 1)
	s.SetID(atomic.AddUint64(&m.serial, 1))
	m.alive.Store(s.ID(), s)
}

func (m *sessionManager[S]) removeAlive(id uint64) {
	atomic.AddInt32(&m.aliveNum, -1)
	m.alive.Delete(id)
}

func (m *sessionManager[S]) getAlive(id uint64) (ISession, bool) {
	v, ok := m.alive.Load(id)
	if !ok {
		return nil, false
	}
	return v.(ISession), true
}

func (m *sessionManager[S]) count() int32 {
	return atomic.LoadInt32(&m.aliveNum)
}

func (m *sessionManager[S]) getIdle() S {
	return m.pool.Get().(S)
}

func (m *sessionManager[S]) putIdle(s S) {
	m.pool.Put(s)
}
