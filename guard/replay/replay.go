package replay

import (
	"net/http"
	"time"

	"github.com/silentred/glog"

	"fmt"

	"github.com/silentred/gateway/reactor/circuit_breaker"
	"github.com/silentred/gateway/util"
)

var (
	ErrDupNonce = util.NewError(100410, "dupliacted nonce")
)

type ReplayReactor struct {
	kv             circuit_breaker.TTLKV
	EffectDuration time.Duration
}

func NewReplayReactor(kv circuit_breaker.TTLKV, effect time.Duration) *ReplayReactor {
	return &ReplayReactor{
		kv:             kv,
		EffectDuration: effect,
	}
}

func (rr *ReplayReactor) Reject(r *http.Request) error {
	var nonce string
	var err, returnErr error
	var val int

	if r != nil {
		nonce = r.Header.Get("X-Nonce")
		if len(nonce) > 0 {
			val, err = rr.kv.Get(nonce)
			// not found
			if err == nil || (err != circuit_breaker.ErrKeyNotFound && val > 0) {
				returnErr = ErrDupNonce
			}

			// save nonce to kv
			err = rr.kv.Set(nonce, 1, rr.EffectDuration)
			if err != nil {
				glog.Errorf("[replay] err=%v", err)
			}
		}
	}

	return returnErr
}

func (rr *ReplayReactor) String() string {
	return fmt.Sprintf("replay(%s)", rr.EffectDuration)
}
