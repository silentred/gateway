package metric

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/silentred/glog"
	"github.com/silentred/gateway/route"
)

const (
	HTTPConnection      = "http"
	WebsocketConnection = "ws"

	IncreaseOpt = 1
	DecreaseOpt = 2
)

var (
	// httpDuration = prometheus.NewSummaryVec(
	// 	prometheus.SummaryOpts{
	// 		Namespace:  "iget",
	// 		Subsystem:  "gw",
	// 		Name:       "http_durations_us",
	// 		Help:       "http latency distributions.",
	// 		Objectives: prometheus.DefObjectives,
	// 	},
	// 	[]string{"svc", "path"},
	// )

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "iget",
			Subsystem: "entree",
			Name:      "http_latency_second",
			Help:      "http latency distributions.",
			// 10 buckets from 20ms to 520ms
			Buckets: prometheus.LinearBuckets(0.02, 0.05, 10),
		},
		[]string{"svc", "path", "code"},
	)

	connectionNum = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "iget",
			Subsystem: "entree",
			Name:      "connection_number",
			Help:      "connection number at the moment",
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(httpDuration, connectionNum)
}

// Start collecting metrics
func Start(path, listenAddr string) {
	var err error
	http.Handle(path, promhttp.Handler())
	http.HandleFunc("/health", Health)
	err = http.ListenAndServe(listenAddr, nil)
	glog.Errorf("[Metric] err=%v", err)
}

// HTTPDurationObserve observes http response
func HTTPDurationObserve(resp *http.Response, svc *route.Service, d time.Duration) {
	if resp != nil && svc != nil {
		httpDuration.WithLabelValues(svc.Name, resp.Request.URL.Path, strconv.Itoa(resp.StatusCode)).Observe(d.Seconds())
	}
}

// AddConnectionNumber adds number by delta
func AddConnectionNumber(delta int, typeName string) {
	connectionNum.WithLabelValues(typeName).Add(float64(delta))
}
