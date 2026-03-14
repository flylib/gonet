package gonet

import "sync"

type bytePool struct {
	pool   sync.Pool
	maxCap int
}

func newBytePool(maxCap int) bytePool {
	return bytePool{maxCap: maxCap}
}

func (p *bytePool) get(size int) []byte {
	if size <= 0 {
		return nil
	}
	if p.maxCap <= 0 {
		return make([]byte, size)
	}
	if v := p.pool.Get(); v != nil {
		buf := v.([]byte)
		if cap(buf) >= size {
			return buf[:size]
		}
	}
	return make([]byte, size)
}

func (p *bytePool) put(buf []byte) {
	if p.maxCap <= 0 || buf == nil {
		return
	}
	if cap(buf) > p.maxCap {
		return
	}
	p.pool.Put(buf[:0])
}
