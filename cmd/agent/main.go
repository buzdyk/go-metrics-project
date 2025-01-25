package main

import (
	"fmt"
	"github.com/buzdyk/go-metrics-project/internal/metrics"
	"net/http"
	"time"
)

func main() {
	for {
		time.Sleep(2 * time.Second)
		go func() {
			r, err := http.Post("http://127.0.0.1:8080/update/gauge/alloc/"+string(metrics.Alloc()), "text/plain", nil)

			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(r.StatusCode)
			fmt.Println(metrics.Alloc())
		}()
	}
}
