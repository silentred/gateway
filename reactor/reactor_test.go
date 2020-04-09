package reactor

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReactor(t *testing.T) {
	g := NewGroup("test", &MockReactor{})
	req, _ := http.NewRequest(http.MethodGet, "http://www.baidu.com/v1/A/hello?a=b&c=d", nil)
	err := g.Reject(req)
	assert.NoError(t, err)
	sign := g.String()
	t.Log(sign)

	var resp http.Response
	err = g.Modify(&resp)
	assert.NoError(t, err)

	g.Add(&MockReactor{})

}
