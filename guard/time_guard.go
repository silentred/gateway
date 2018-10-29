package guard

import (
	"net/http"
	"strconv"
	"time"

	"github.com/silentred/glog"
	"github.com/silentred/gateway/util"
)

const (
	CodeTimestamp = 100405
)

var (
	ErrTimestamp = util.NewError(CodeTimestamp, "invalid timestamp")
)

// TimeGuard check the correctness of timestamp in header
type TimeGuard struct {
	now   func() time.Time
	delta int64
}

// NewTimeGuard returns a TimeGuard
func NewTimeGuard(now func() time.Time, delta int64) *TimeGuard {
	return &TimeGuard{
		now:   now,
		delta: delta,
	}
}

// Reject implements Guard interface
func (tg *TimeGuard) Reject(r *http.Request) error {
	var timeStr string
	var timestamp int
	var now int64
	var err error

	timeStr = r.Header.Get("X-Timestamp")
	timestamp, err = strconv.Atoi(timeStr)
	if err != nil {
		return err
	}

	now = tg.now().Unix()
	glog.Debugf("[Time] time=%d now=%d", timestamp, now)

	if int64(timestamp) < now-tg.delta || int64(timestamp) > now+tg.delta {
		return ErrTimestamp.Now()
	}

	return nil
}

func (tg *TimeGuard) String() string {
	return "time-guard"
}
