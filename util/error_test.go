package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := NewError(1, "test")
	assert.Equal(t, 1, err.Code)
	assert.Equal(t, "test", err.Msg)
	t.Log(err.Error())
	err.Timestamp = 123
	t.Log(err.Error())
}
