package proxy

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/silentred/glog"
	"github.com/silentred/gateway/util"
)

const (
	CodeTimeoutErr = 100501
	CodeUnkownErr  = 100600
)

type ErrWriter struct {
	rw http.ResponseWriter
}

func NewErrWriter(rw http.ResponseWriter) *ErrWriter {
	return &ErrWriter{rw}
}

func (ew *ErrWriter) Write(p []byte) (int, error) {
	var err util.Error
	var e error
	var b []byte
	if bytes.Contains(p, []byte("timeout")) {
		err = util.NewError(CodeTimeoutErr, util.String(p))
	} else {
		err = util.NewError(CodeUnkownErr, util.String(p))
	}
	b, e = json.Marshal(err)
	if e != nil {
		glog.Errorf("[JSON] err=%v", err)
	}
	return ew.rw.Write(b)
}
