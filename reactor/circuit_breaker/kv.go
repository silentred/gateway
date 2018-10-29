package circuit_breaker

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/silentred/gateway/util"
)

const (
	CodeKeyNotFound = 200100
)

var (
	ErrKeyNotFound = util.NewError(CodeKeyNotFound, "key not found")
)

type TTLKV interface {
	Get(key string) (int, error)
	Set(key string, val int, ttl time.Duration) error
}

type MemTTLKV struct {
	TTL   time.Duration
	Cache *cache.Cache
}

func NewMemTTLKV(ttl time.Duration) *MemTTLKV {
	c := cache.New(ttl, ttl/2)
	return &MemTTLKV{
		TTL:   ttl,
		Cache: c,
	}
}

func (mk *MemTTLKV) Set(key string, val int, ttl time.Duration) error {
	mk.Cache.Set(key, val, ttl)
	return nil
}

func (mk *MemTTLKV) Get(key string) (int, error) {
	val, found := mk.Cache.Get(key)
	if i, ok := val.(int); found && ok {
		return i, nil
	}
	return 0, ErrKeyNotFound
}
