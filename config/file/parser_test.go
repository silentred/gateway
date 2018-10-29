package file

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/silentred/gateway/route"
)

func TestParse(t *testing.T) {
	line := `Route:www.luojilab.com Prefix:/v1/A/hello Service:serviceA Strip:/v1/A Targets(127.0.0.1:8080,100)`
	route, svc, err := parse(line)
	assert.NoError(t, err)
	assert.Equal(t, "www.luojilab.com", route.host)
	assert.Equal(t, "serviceA", svc.name)
}

func TestFileBackend(t *testing.T) {
	table := route.NewTable()
	fb := NewFileBackend("../../route.cfg", table)
	err := fb.Parse()
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodGet, "http://www.luojilab.com/v1/A/hello?a=b&c=d", nil)

	r := route.NewRoute("www.luojilab.com", "/v1/A/hello")
	s := table.Find(r)
	assert.NotNil(t, s)
	assert.Equal(t, "serviceA", s.Name)

	s = table.FindByRequest(req)
	assert.NotNil(t, s)
	assert.Equal(t, "serviceA", s.Name)

	r = route.NewRoute("s2.luojilab.com", "/v1/B/world")
	s = table.Find(r)
	assert.NotNil(t, s)
	assert.Equal(t, "serviceB", s.Name)

	go fb.Watch()
	// f, err := os.OpenFile("../route.cfg", os.O_APPEND|os.O_WRONLY, 0755)
	// assert.NoError(t, err)
	// f.WriteString("\nRoute:s3.luojilab.com Prefix:/v1/C/hello Service:serviceC Strip:/v1/C Targets(127.0.0.1:8082,100)")
	// f.Sync()
	// f.Close()

	time.Sleep(1 * time.Second)
}
