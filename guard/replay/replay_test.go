package replay

import (
	"net/http"
	"testing"
	"time"

	"github.com/silentred/kassadin/util"

	"github.com/stretchr/testify/assert"
	"github.com/silentred/gateway/reactor/circuit_breaker"
)

func TestReplay(t *testing.T) {
	var err error

	c := circuit_breaker.NewMemTTLKV(2 * time.Second)
	rr := NewReplayReactor(c, time.Second)

	header := map[string]string{
		"X-Nonce": "randomstr",
	}
	req, _ := util.NewHTTPReqeust(http.MethodGet, "http://www.luojilab.com/sdf", nil, header, nil)

	// first visit, pass
	err = rr.Reject(req)
	assert.NoError(t, err)

	// second visit, found existing nonce
	err = rr.Reject(req)
	assert.Equal(t, ErrDupNonce, err)

	time.Sleep(time.Second)

	// after 1s, pass again
	err = rr.Reject(req)
	assert.NoError(t, err)
}
