package route

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	var err error
	rt := NewRoute("www.fakebaidu.com", "/v1/hello")
	target := NewTarget("", "www.baidu.com", 1)
	service := NewService("test service", "", "/v1", target, nil, nil)

	tbl := NewTable()
	err = tbl.Reset()
	assert.NoError(t, err)

	err = tbl.Add(rt, service)
	assert.NoError(t, err)

	svc := tbl.Find(rt)
	assert.Equal(t, "test service", svc.Name)

	ok := tbl.Del(&rt)
	assert.True(t, ok)

	err = tbl.Add(rt, service)
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "http://www.fakebaidu.com/v1/hello", nil)
	svc = tbl.FindByRequest(req)
	assert.Equal(t, "test service", svc.Name)
}
