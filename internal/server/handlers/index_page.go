package handlers

import (
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"
)

func (mh *MetricHandler) GetIndex(rw http.ResponseWriter, r *http.Request) {
	gauges, _ := mh.gaugeStore.Values(r.Context())
	counters, _ := mh.counterStore.Values(r.Context())

	data := struct {
		Gauges   map[string]metrics.Gauge
		Counters map[string]metrics.Counter
	}{gauges, counters}

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	tmpl, err := template.ParseFiles(dir + "/templates/index.html")
	if err != nil {
		http.Error(rw, "Failed to parse HTML template", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(rw, data); err != nil {
		fmt.Println(err)
		http.Error(rw, "Failed to render metrics page", http.StatusInternalServerError)
	}
}
