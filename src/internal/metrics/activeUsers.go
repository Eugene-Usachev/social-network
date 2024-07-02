package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Active Users Gauge
var activeUsers = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "social_network",
	Name:      "active_users",
	Help:      "Current number of active users",
})

// IncActiveUsers increments Active Users Gauge
func IncActiveUsers() {
	activeUsers.Add(float64(1))
}

// DecActiveUsers decrements Active Users Gauge
func DecActiveUsers() {
	activeUsers.Add(float64(-1))
}
