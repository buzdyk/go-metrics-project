package main

import (
	"net/http"
	"time"
)

func main() {
	for {
		time.Sleep(2 * time.Second)
		http.Get(":8080")
	}
}
