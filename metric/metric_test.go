package metric

import (
	"net/http"
	"testing"
	"time"

	"github.com/silentred/gateway/route"
)

func TestMetrics(t *testing.T) {
	go Start("metrics", ":7088")

	req, _ := http.NewRequest(http.MethodGet, "http://www.fakeweb.com/v1/hello?a=b&c=d", nil)
	resp := &http.Response{Request: req}
	svc := &route.Service{Name: "test"}
	HTTPDurationObserve(resp, svc, time.Second)
}
