package circuit_breaker

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/silentred/glog"
	"github.com/silentred/gateway/route"
	"github.com/silentred/gateway/util"
)

const (
	CodeCircuitBreak = 200200
)

var (
	ErrCircuitBreak = util.NewError(CodeCircuitBreak, "circuit break")
)

type CircuitBreaker struct {
	kv       TTLKV
	strategy BlockStrategy
	mut      sync.Mutex

	//
	ThresholdCount    int
	ThresholdDuration time.Duration
	BlockDuration     time.Duration
}

func NewCircuitBreaker(kv TTLKV, s BlockStrategy, tc int, td, bd time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		kv:                kv,
		strategy:          s,
		mut:               sync.Mutex{},
		ThresholdCount:    tc,
		ThresholdDuration: td,
		BlockDuration:     bd,
	}
}

func (cb *CircuitBreaker) Reject(r *http.Request) error {
	var id uint32
	var blockKey string
	id = route.HashID(r.Host, r.URL.Path)
	blockKey = fmt.Sprintf("%d_blk", id)
	_, err := cb.kv.Get(blockKey)

	if err != nil {
		glog.Infof("[CB] try block req host=%s path=%s key=%s kvErr=%v", r.Host, r.URL.Path, blockKey, err)
	}

	if err == ErrKeyNotFound {
		glog.Debugf("decreasing block prob, %s%s", r.Host, r.URL.Path)
		cb.strategy.DecreaseProbability(id)
	} else if cb.strategy.Block(id) {
		glog.Infof("[CB] blockd req host=%s path=%s", r.Host, r.URL.Path)
		return ErrCircuitBreak
	}

	return nil
}

// Modify cannot observe request, because Director() has changed the request's host and path.
// This method is only capable to observe response.
func (cb *CircuitBreaker) Modify(resp *http.Response) error {
	// var id uint32
	// var currCount int

	// id = route.HashID(resp.Request.Host, resp.Request.URL.Path)
	// // resp error: increase count, rare operation
	// if resp != nil && resp.StatusCode >= http.StatusInternalServerError {
	// 	glog.Debugf("record bad route-id=%d host=%s path=%s", id, resp.Request.Host, resp.Request.URL.Path)
	// 	currCount = cb.recordError(id)
	// }
	// // resp is ok: check count, set block flag, add block probability
	// cb.increaseBlockProb(id, currCount)
	return nil
}

// ObserveError should be used to observe original request and resonse error.
// The request of first parameter should be the original request.
func (cb *CircuitBreaker) ObserveError(req *http.Request, resp *http.Response, err error) {
	var id uint32
	var currCount int

	id = route.HashID(req.Host, req.URL.Path)
	glog.Debugf("record bad route-id=%d host=%s path=%s", id, req.Host, req.URL.Path)
	if err != nil || resp.StatusCode >= http.StatusInternalServerError {
		currCount = cb.recordError(id)
	}
	// resp is ok: check count, set block flag, add block probability
	cb.increaseBlockProb(id, currCount)
}

func (cb *CircuitBreaker) String() string {
	return fmt.Sprintf("circuit-breaker(%d,%s,%s)", cb.ThresholdCount, cb.ThresholdDuration, cb.BlockDuration)
}

func (cb *CircuitBreaker) recordError(routeID uint32) int {
	var currCount int
	var kvErr error
	var key string

	cb.mut.Lock()
	defer cb.mut.Unlock()
	key = fmt.Sprintf("%d", routeID)
	currCount, kvErr = cb.kv.Get(key)
	glog.Debugf("[CB] bad resp key=%s count=%d kvErr=%v", key, currCount, kvErr)
	if kvErr == ErrKeyNotFound {
		cb.kv.Set(strconv.Itoa(int(routeID)), 1, cb.ThresholdDuration)
	} else if kvErr == nil {
		cb.kv.Set(strconv.Itoa(int(routeID)), currCount+1, cb.ThresholdDuration)
	} else {
		glog.Errorf("[CB] err=%v", kvErr)
	}

	return currCount
}

func (cb *CircuitBreaker) increaseBlockProb(routeID uint32, count int) {
	if count > cb.ThresholdCount {
		glog.Debugf("increasing block prob")
		var blockKey = fmt.Sprintf("%d_blk", routeID)
		var err error

		err = cb.kv.Set(blockKey, 0, cb.BlockDuration)
		cb.strategy.IncreaseProbability(routeID)
		if err != nil {
			glog.Errorf("[CB] err=%v", err)
		}

		// reset count of routeID
		var key string
		key = fmt.Sprintf("%d", routeID)
		err = cb.kv.Set(key, 0, cb.ThresholdDuration)
		if err != nil {
			glog.Errorf("[CB] err=%v", err)
		}
	}
}
