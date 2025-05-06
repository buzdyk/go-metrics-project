package handlers

import (
	"github.com/buzdyk/go-metrics-project/internal/database"
	_ "github.com/lib/pq"
	"net/http"
)

func (mh *MetricHandler) Ping(rw http.ResponseWriter, r *http.Request) {
	dbClient := database.GetClient()

	if err := dbClient.Ping(); err != nil {
		http.Error(rw, "database ping failed", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Connection ok"))
}
