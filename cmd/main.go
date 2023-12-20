package main

import (
	"metrostreams/internal/pkg/api"
	"net/http"
)

func main() {
	http.HandleFunc("/calculate-stream/", api.Calculate)
	http.ListenAndServe("localhost:8080", nil)
}
