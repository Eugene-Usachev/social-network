package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// User Registrations Counter
var userRegistrations = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "social_network",
	Name:      "user_registrations_total",
	Help:      "Total number of user registrations",
})

// IncUserRegistrations increments User Registrations Counter
func IncUserRegistrations() {
	userRegistrations.Inc()
}
