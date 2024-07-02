package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"social-network/src/pkg/logger"
)

func MustListen(addr string, logger logger.Logger) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start metrics server: %s", err))
	}
}
