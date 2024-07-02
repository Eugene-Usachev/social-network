package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"time"
)

// HTTP requests summary (duration in seconds, method, path, status code)
var httpRequestsMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace: "social_network",
	Name:      "http_requests",
	Help:      "http requests summary",
	Objectives: map[float64]float64{
		0.5:  0.05,
		0.9:  0.01,
		0.99: 0.001,
	},
}, []string{"method", "path", "status_code"})

// ObserveRequest reports an HTTP request summary (duration in seconds, method, path, status code)
func ObserveRequest(duration time.Duration, method string, path string, status int) {
	httpRequestsMetrics.WithLabelValues(method, path, strconv.Itoa(status)).Observe(duration.Seconds())
}
