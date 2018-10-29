package circuit_breaker

import (
	"sync"
	"time"

	"github.com/silentred/glog"
)

type BlockStrategy interface {
	IncreaseProbability(reqID uint32)
	DecreaseProbability(reqID uint32)
	Block(reqID uint32) bool
}

type BinaryStrategy struct {
	// map[reqID]count, every request should have its own factor
	factors map[uint32]uint32
	mut     sync.Mutex
}

func NewBinaryStrategy() *BinaryStrategy {
	return &BinaryStrategy{
		factors: make(map[uint32]uint32),
		mut:     sync.Mutex{},
	}
}

func (bs *BinaryStrategy) IncreaseProbability(reqID uint32) {
	bs.mut.Lock()
	defer bs.mut.Unlock()
	if f, has := bs.factors[reqID]; has {
		bs.factors[reqID] = f * 2
	} else {
		bs.factors[reqID] = 2
	}
}

func (bs *BinaryStrategy) DecreaseProbability(reqID uint32) {
	if f, has := bs.factors[reqID]; has {
		bs.mut.Lock()
		defer bs.mut.Unlock()
		bs.factors[reqID] = f / 2
		glog.Debugf("[CB] f=%d nowF=%d", f, f/2)
		if bs.factors[reqID] < 2 {
			delete(bs.factors, reqID)
		}
	}
	//^uint32(C-1) == -C
	//atomic.AddUint32(&bs.factor, ^uint32(bs.factor/2-1))
}

func (bs *BinaryStrategy) Block(reqID uint32) bool {
	if f, has := bs.factors[reqID]; has {
		if uint64(time.Now().UnixNano())%uint64(f) > 0 {
			return true
		}
	}
	return false
}
