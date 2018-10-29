package guard

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/silentred/glog"
	"github.com/silentred/gateway/util"
)

const (
	CodeInvalidParam = 100401
	CodeSign         = 100403
)

var (
	ErrInvalidParam = util.NewError(CodeInvalidParam, "invalid parameters")
	ErrSign         = util.NewError(CodeSign, "invalid sign")
)

// SignGuard check the correctness of sign
type SignGuard struct {
	secret string
}

// NewSignGuard returns a new SignGuard
func NewSignGuard(secret string) *SignGuard {
	return &SignGuard{secret}
}

// Reject implements the Guard interface
func (sg *SignGuard) Reject(r *http.Request) error {
	var appKey, method, conType, path, queryStr, timestamp, nonce string
	var hashStr, sign, rightSign string

	appKey = r.Header.Get("X-App-Key")
	method = r.Method
	conType = r.Header.Get("Content-Type")
	path = r.URL.Path
	queryStr = r.URL.Query().Encode()
	timestamp = r.Header.Get("X-Timestamp")
	nonce = r.Header.Get("X-Nonce")

	sign = r.Header.Get("X-Sign")
	// not empty
	if len(appKey) == 0 || len(sign) == 0 {
		glog.Debugf("[Sign] appKey=%s conType=%s sign=%s", appKey, conType, sign)
		return ErrInvalidParam
	}

	hashStr = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s", appKey, method, conType, path, queryStr, timestamp, nonce, sg.secret)
	bytes := md5.Sum(util.Slice(hashStr))
	md5Str := fmt.Sprintf("%x", util.String(bytes[:]))
	rightSign = base64.StdEncoding.EncodeToString(util.Slice(md5Str))
	glog.Debugf("[Sign] str=%s sign=%s compareTo=%s", hashStr, rightSign, sign)

	if sign != rightSign {
		return ErrSign
	}

	return nil
}

func (sg *SignGuard) String() string {
	return "sign-guard"
}

func sign(key, method, conType, path, query, time, nonce, secret string) string {
	hashStr := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s", key, method, conType, path, query, time, nonce, secret)
	bytes := md5.Sum(util.Slice(hashStr))
	md5Str := fmt.Sprintf("%x", util.String(bytes[:]))
	return base64.StdEncoding.EncodeToString(util.Slice(md5Str))
}
