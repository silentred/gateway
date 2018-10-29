package rate_limiter

import (
	"fmt"
	"net/http"

	"github.com/silentred/gateway/util"

	"github.com/silentred/glog"
)

var (
	ErrRateLimit = util.NewError(100420, "reached rate limit")
)

type RateLimiter struct {
	kv           AtomicTTLKV
	ThresholdCnt int
}

func NewRateLimiter(kv AtomicTTLKV, cnt int) *RateLimiter {
	return &RateLimiter{
		kv:           kv,
		ThresholdCnt: cnt,
	}
}

func (rl *RateLimiter) Reject(req *http.Request) error {
	var uidStr, uidKey string
	var count int
	var err error

	uidStr = req.Header.Get("X-UID")
	if len(uidStr) == 0 {
		glog.Warningf("[RateLimiter] uidKey is empty")
		return nil
	}

	uidKey = fmt.Sprintf("rl:%s", uidStr)

	count, err = rl.kv.Get(uidKey)
	if err != nil {
		glog.Errorf("[RateLimiter] uidKey=%s err=%v cnt=%d", uidKey, err, count)
	}

	if count > rl.ThresholdCnt {
		return ErrRateLimit
	}

	err = rl.kv.Incr(uidKey, 1)
	if err != nil {
		glog.Errorf("[RateLimiter] uidKey=%s err=%v", uidKey, err)
	}

	return nil
}

func (rl *RateLimiter) String() string {
	return fmt.Sprintf("rate_limiter(%d)", rl.ThresholdCnt)
}
