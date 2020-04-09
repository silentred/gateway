package circuit_breaker

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker(t *testing.T) {
	c := NewMemTTLKV(3 * time.Second)
	s := NewBinaryStrategy()
	var thresholdCnt = 10
	var thresholdDuration = time.Second
	var blockDuration = time.Second
	var err error
	cb := NewCircuitBreaker(c, s, thresholdCnt, thresholdDuration, blockDuration)

	req, _ := http.NewRequest(http.MethodGet, "http://www.baidu.com/v1/A/hello?a=b&c=d", nil)
	resp := &http.Response{Request: req, StatusCode: 500}

	// not reject, because not reaching the thresholdCnt
	for i := 0; i < 10; i++ {
		cb.ObserveError(req, resp, nil)
		assert.NoError(t, err)
	}
	err = cb.Reject(req)
	assert.NoError(t, err)

	// reject
	var rejectCnt, passCnt int
	for i := 0; i < 20; i++ {
		cb.ObserveError(req, resp, nil)
		err = cb.Reject(req)
		if err == ErrCircuitBreak {
			rejectCnt++
		} else {
			passCnt++
		}
	}
	assert.True(t, rejectCnt > 0)
	t.Logf("pass=%d reject=%d", passCnt, rejectCnt)

	// not reject, because block duration has passed
	time.Sleep(time.Second)
	err = cb.Reject(req)
	assert.NoError(t, err)
}
