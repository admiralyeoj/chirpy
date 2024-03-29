package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

// headers is a handler function for the metrics endpoint
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	// Set the response headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Respond with a 200 OK status and "OK" message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	</html>`, cfg.fileserverHits)))
}

// headers is a handler function for the reset endpoint
func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	fmt.Printf("cfg.fileserverHits = %d", cfg.fileserverHits)

	// Respond with a 200 OK status and "OK" message
	w.WriteHeader(http.StatusOK)
}
