package handlers

import (
	"database/sql"
	"github.com/buzdyk/go-metrics-project/internal/server/config"
	_ "github.com/lib/pq"
	"net/http"
)

func (mh *MetricHandler) Ping(rw http.ResponseWriter, r *http.Request) {
	dsn := config.GetConfig().PgDsn

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		http.Error(rw, "failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		http.Error(rw, "database ping failed", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Connection ok"))
}
