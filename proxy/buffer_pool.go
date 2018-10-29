package proxy

import (
	"sync"
)

type bufferPool struct {
	size int
	pool sync.Pool
}

func NewBufferPool(size int) *bufferPool {
	bp := &bufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, size)
			},
		},
	}
	return bp
}

func (bp *bufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

func (bp *bufferPool) Put(b []byte) {
	bp.pool.Put(b)
}
