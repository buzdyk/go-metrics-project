package server

import (
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"github.com/buzdyk/go-metrics-project/internal/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var store = storage.NewMemStorage()

var StoreMetric = func(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")
	metricValue := r.PathValue("value")

	//if metrics.Exists(r.PathValue("metric")) == false {
	//	http.Error(rw, "metric does not exist", http.StatusBadRequest)
	//	return
	//}

	switch metricType {
	case "gauge":
		v, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "metric value is not convertible to float64", http.StatusBadRequest)
		}
		store.StoreGauge(metricName, metrics.Gauge(v))
	case "counter":
		v, err := strconv.Atoi(metricValue)
		if err != nil {
			http.Error(rw, "metric value is not convertible to int", http.StatusBadRequest)
		}
		store.StoreCounter(metricName, metrics.Counter(v))
	}

	log.Default().Println("type:", metricType, "metric", metricName, metricValue)

	rw.WriteHeader(200)
	rw.Write([]byte("ok"))
}

var GetMetric = func(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")

	//if metrics.Exists(r.PathValue("metric")) == false {
	//	http.Error(rw, "metric does not exist", http.StatusBadRequest)
	//	return
	//}

	switch metricType {
	case "gauge":
		if v, err := store.Gauge(metricName); err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 64)))
		}
	case "counter":
		v, err := store.Counter(metricName)
		if err != nil {
			rw.WriteHeader(404)
		} else {
			rw.Write([]byte(strconv.Itoa(int(v))))
		}
	}
}

var GetIndex = func(rw http.ResponseWriter, r *http.Request) {
	const templateString = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Metrics Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; padding: 20px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <h1>Metrics Dashboard</h1>

    <h2>Gauges</h2>
    <table>
        <tr><th>Name</th><th>Value</th></tr>
        {{range $key, $value := .Gauges}}
        <tr><td>{{$key}}</td><td>{{$value}}</td></tr>
        {{end}}
    </table>

    <h2>Counters</h2>
    <table>
        <tr><th>Name</th><th>Value</th></tr>
        {{range $key, $value := .Counters}}
        <tr><td>{{$key}}</td><td>{{$value}}</td></tr>
        {{end}}
    </table>
</body>
</html>
`

	data := struct {
		Gauges   map[string]metrics.Gauge
		Counters map[string]metrics.Counter
	}{
		store.Gauges(),
		store.Counters(),
	}

	tmpl, err := template.New("metrics").Parse(templateString)
	if err != nil {
		http.Error(rw, "Failed to parse HTML template", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(rw, data); err != nil {
		http.Error(rw, "Failed to render metrics page", http.StatusInternalServerError)
	}

}
