package util

import (
	"encoding/json"
	"time"
)

// Error implements error
type Error struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"ts,omitempty"`
}

// NewError returns a new Error
func NewError(code int, msg string) Error {
	return Error{
		Code: code,
		Msg:  msg,
	}
}

func (e Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// Now set the Timestamp to now
func (e Error) Now() Error {
	e.Timestamp = time.Now().Unix()
	return e
}
