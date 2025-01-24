package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type gauge float64
type counter int64

type MemStorage struct {
	g map[string][]gauge
	c map[string]counter
}

func (s *MemStorage) StoreGauge(name string, value gauge) {
	s.g[name] = append(s.g[name], value)
}

func (s *MemStorage) AddCounter(name string, value counter) {
	if _, ok := s.c[name]; !ok {
		s.c[name] = 0
	}
	s.c[name] += value
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		g: make(map[string][]gauge),
		c: make(map[string]counter),
	}
}

var storage *MemStorage

func init() {
	storage = NewMemStorage()
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /update/{type}/{metric}/{value}", storeMetric)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func storeMetric(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")
	metricValue, _ := strconv.Atoi(r.PathValue("value"))

	switch metricType {
	case "gauge":
		storage.StoreGauge(metricName, gauge(metricValue))
	case "counter":
		storage.AddCounter(metricName, counter(metricValue))
	}

	log.Default().Println("type:", metricType, "metric", metricName)

	rw.WriteHeader(200)
	rw.Write([]byte("ok"))

	fmt.Println(storage)
}
