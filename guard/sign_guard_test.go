package guard

import (
	"net/http"
	"net/http/httputil"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/silentred/gateway/util"
)

func TestSignGuard(t *testing.T) {
	secret := "abc123"
	key := "ios-client"
	g := NewSignGuard(secret)
	query := map[string]string{
		"a": "1",
		"b": "2",
	}
	time := strconv.Itoa(int(time.Now().Unix()))
	header := map[string]string{
		//"X-App-Key":    key,
		"Content-Type": "application/json",
		"X-Timestamp":  time,
		"X-Nonce":      "123",
		"X-Sign":       "wrong sign",
	}
	// for invalid param
	req, err := util.NewHTTPReqeust(http.MethodGet, "http://www.luojilab.com/sdf", query, header, nil)
	assert.NoError(t, err)
	err = g.Reject(req)
	assert.Equal(t, ErrInvalidParam, err)

	// for invalid sign
	header["X-App-Key"] = key
	req, err = util.NewHTTPReqeust(http.MethodGet, "http://www.luojilab.com/sdf", query, header, nil)
	assert.NoError(t, err)
	err = g.Reject(req)
	assert.Equal(t, ErrSign, err)

	// for right sign
	s := sign(key, http.MethodGet, "application/json", req.URL.Path, req.URL.Query().Encode(), time, "123", secret)
	header["X-Sign"] = s
	req, err = util.NewHTTPReqeust(http.MethodGet, "http://www.luojilab.com/sdf", query, header, nil)
	b, _ := httputil.DumpRequest(req, false)
	t.Log(string(b))
	assert.NoError(t, err)
	err = g.Reject(req)
	assert.NoError(t, err)

	// test
	// s = sign("iget-ios-client-key", "POST", "application/json; charset=UTF-8", "/v3/bookrack/list", "", "1494591431", "123", "test123")
	// t.Log(s)
}
