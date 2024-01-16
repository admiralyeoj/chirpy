package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "./app"
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	// Create a new http.ServeMux
	mux := http.NewServeMux()

	// Wrap the mux with the middlewareCors
	corsMux := middlewareCors(mux)

	// Create a handler for serving static files from the "/app/" path
	// Change the path to your actual directory where index.html and logo.png are located
	appDir := http.Dir(filepathRoot)
	appHandler := http.StripPrefix("/app", http.FileServer(appDir))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(appHandler))

	// Register the healthz endpoint using mux.HandleFunc
	mux.HandleFunc("/healthz", handlerReadiness)
	// Register the hits endpoint using mux.HandleFunc
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	// Register the hits endpoint using mux.HandleFunc
	mux.HandleFunc("/reset", apiCfg.handlerReset)

	// Create a new http.Server and use corsMux as the handler
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
