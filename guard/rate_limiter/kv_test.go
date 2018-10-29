package rate_limiter

import (
	"testing"
	"time"

	"sync"

	"github.com/stretchr/testify/assert"
	"github.com/silentred/gateway/reactor/circuit_breaker"
)

func TestAtomicTTLKV(t *testing.T) {
	var err error
	key := "test123"

	attlKV := newAtomicTTLKV()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			err = attlKV.Incr(key, 1)
			assert.NoError(t, err)
			wg.Done()
		}()
	}
	wg.Wait()

	val, err := attlKV.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, 10, val)
}

func newAtomicTTLKV() *MemAtomicTTLKV {
	ttlKV := circuit_breaker.NewMemTTLKV(time.Second)
	attlKV := NewMemAtomicTTLKV(*ttlKV, 10)
	return attlKV
}
