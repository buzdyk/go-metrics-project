package server

import (
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"net/http"
)

func metricExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := metrics.Collectors[r.PathValue("metric")]; ok {
			http.Error(w, "metric does not exist", http.StatusBadRequest)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
