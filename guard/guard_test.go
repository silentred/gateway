package guard

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuard(t *testing.T) {
	g := NewGroup("test", &MockGuard{})
	req, _ := http.NewRequest(http.MethodGet, "http://www.baidu.com/v1/A/hello?a=b&c=d", nil)
	err := g.Reject(req)
	assert.NoError(t, err)
	t.Log(g.String())
	g.Add(&MockGuard{})
}
