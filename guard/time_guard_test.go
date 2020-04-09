package guard

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/silentred/gateway/util"
	"github.com/stretchr/testify/assert"
)

func TestTimeGuard(t *testing.T) {
	now := func() time.Time {
		return time.Now()
	}

	g := NewTimeGuard(now, 15*60)

	timeStr := strconv.Itoa(int(time.Now().Unix()))
	header := map[string]string{
		"X-Timestamp": timeStr,
	}
	// for no error
	req, err := util.NewHTTPReqeust(http.MethodGet, "http://www.baidu.com/sdf", nil, header, nil)
	assert.NoError(t, err)
	err = g.Reject(req)
	assert.NoError(t, err)

	// for invalid time
	past := time.Now().Add(-30 * time.Minute).Unix()
	header["X-Timestamp"] = strconv.Itoa(int(past))
	req, _ = util.NewHTTPReqeust(http.MethodGet, "http://www.baidu.com/sdf", nil, header, nil)
	err = g.Reject(req)
	assert.Equal(t, ErrTimestamp.Now(), err)
}
