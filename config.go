package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

// headers is a handler function for the metrics endpoint
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	// Set the response headers
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Respond with a 200 OK status and "OK" message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + strconv.Itoa(cfg.fileserverHits)))
}

// headers is a handler function for the reset endpoint
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	fmt.Println("cfg.fileserverHits = " + strconv.Itoa(cfg.fileserverHits))

	// Respond with a 200 OK status and "OK" message
	w.WriteHeader(http.StatusOK)
}
