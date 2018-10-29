package circuit_breaker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTTLKV(t *testing.T) {
	c := NewMemTTLKV(time.Second)
	err := c.Set("foo", 10, time.Second/2)
	assert.NoError(t, err)
	val, err := c.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, 10, val)

	time.Sleep(time.Second / 2)
	_, err = c.Get("foo")
	assert.Equal(t, ErrKeyNotFound, err)

}
