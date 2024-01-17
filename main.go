package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	const filepathRoot = "./app"
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	// Create a new http.ServeMux
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Wrap the mux with the middlewareCors
	corsMux := middlewareCors(r)

	// Create a handler for serving static files from the "/app/" path
	// Change the path to your actual directory where index.html and logo.png are located
	// appDir := http.Dir(filepathRoot)
	// appHandler := http.StripPrefix("/app", http.FileServer(appDir))
	// r.Handle("/app/", apiCfg.middlewareMetricsInc(appHandler))

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)

	// Register the file server to handle requests for assets.
	r.Get("/app/*", func(w http.ResponseWriter, r *http.Request) {
		fsHandler.ServeHTTP(w, r)
	})

	r.Get("/healthz", handlerReadiness)
	r.Get("/metrics", apiCfg.handlerMetrics)

	// Register the hits endpoint using mux.HandleFunc
	r.HandleFunc("/reset", apiCfg.handlerReset)

	// Create a new http.Server and use corsMux as the handler
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
