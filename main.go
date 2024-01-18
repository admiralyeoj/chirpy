package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/admiralyeoj/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	if *dbg {
		fmt.Println("Debug mode is enabled")
		err := os.Remove(dbPath)
		if err != nil {
			fmt.Println("Error deleting the file:", err)
			return
		}
	}

	const filepathRoot = "./app"
	const port = "8080"

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewareCors)

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	adminRouter := chi.NewRouter()
	r.Mount("/admin", adminRouter)

	adminRouter.Get("/metrics", apiCfg.handlerMetrics)

	// Mount API Routes Here
	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)

	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerMetricsReset)

	// Chirps
	apiRouter.Post("/chirps", apiCfg.handlerCreateChirp)
	apiRouter.Get("/chirps", apiCfg.handlerGetChirps)
	apiRouter.Get("/chirps/{chirpId}", apiCfg.handlerGetChirpById)

	// Users
	apiRouter.Post("/users", apiCfg.handlerCreateUser)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	http.ListenAndServe(":"+port, r)
}
