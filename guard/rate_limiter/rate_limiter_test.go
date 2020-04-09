package rate_limiter

import (
	"net/http"
	"testing"
	"time"

	"github.com/silentred/kassadin/util"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	var err error

	attlKV := newAtomicTTLKV()
	rl := NewRateLimiter(attlKV, 10)

	header := map[string]string{
		"X-UID": "123",
	}
	req, _ := util.NewHTTPReqeust(http.MethodGet, "http://www.baidu.com/sdf", nil, header, nil)

	for i := 0; i < 11; i++ {
		err = rl.Reject(req)
		assert.NoError(t, err)
	}

	err = rl.Reject(req)
	assert.Equal(t, ErrRateLimit, err)

	time.Sleep(time.Second)

	err = rl.Reject(req)
	assert.NoError(t, err)
}
