package echorus

import (
	"testing"

	"github.com/labstack/gommon/log"
)

func TestPrint(t *testing.T) {
	l := NewLogger()
	l.Debugj(log.JSON{"test": "gogo", "msg": "gogogo "})
	l.Debug("test")
	l.Debugf("aa %s", "gogo")
	l.Infof("%v", []int{1, 2, 3})

	a := log.JSON{"a": "b"}
	b := log.JSON{"c": "d"}
	c := l.MergeJSON(a, b)
	t.Log(c)
}
