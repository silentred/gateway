package route

import "testing"
import "github.com/stretchr/testify/assert"

func TestRoute(t *testing.T) {
	r := NewRoute("www.baidu.com", "/v1/hello")
	assert.Equal(t, "/v1/hello", r.Prefix)

	id := HashID(r.Host, r.Prefix)
	t.Log(id)
	assert.Equal(t, uint32(1637997894), id)
}
