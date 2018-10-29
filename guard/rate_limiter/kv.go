package rate_limiter

import (
	"hash/crc32"
	"sync"

	"github.com/silentred/glog"
	"github.com/silentred/gateway/reactor/circuit_breaker"
	"github.com/silentred/gateway/util"
)

type AtomicTTLKV interface {
	circuit_breaker.TTLKV
	// Increase string key by int n
	Incr(string, int) error
}

type MemAtomicTTLKV struct {
	circuit_breaker.MemTTLKV
	lockNum int
	muts    []sync.Mutex
}

func NewMemAtomicTTLKV(kv circuit_breaker.MemTTLKV, lockNum int) *MemAtomicTTLKV {
	muts := make([]sync.Mutex, lockNum)
	return &MemAtomicTTLKV{
		MemTTLKV: kv,
		lockNum:  lockNum,
		muts:     muts,
	}
}

func (mkv *MemAtomicTTLKV) Incr(key string, delta int) error {
	var index uint32
	var val int
	var err error

	index = crc32.ChecksumIEEE(util.Slice(key)) % uint32(len(mkv.muts))
	mkv.muts[index].Lock()
	defer mkv.muts[index].Unlock()

	val, err = mkv.Get(key)
	if err != nil {
		glog.Errorf("[AtomicTTLKV] err=%v", err)
	}
	val += delta
	err = mkv.Set(key, val, mkv.TTL)
	if err != nil {
		glog.Errorf("[AtomicTTLKV] err=%v", err)
	}

	return nil
}
