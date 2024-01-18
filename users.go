package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/admiralyeoj/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Create a variable to hold the JSON data.
	var body RequestBody

	// Decode the JSON data from the request body.
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	chirp, err := cfg.DB.CreateUser(body.Email)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	// Create a variable to hold the JSON data.
	db, err := database.NewDB(dbPath)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	chirps, err := db.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetUserById(w http.ResponseWriter, r *http.Request) {
	// Create a variable to hold the JSON data.
	chirpIDStr := chi.URLParam(r, "chirpId")
	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		// Handle the error if the parameter is not a valid integer
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirps, err := cfg.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	chirp := database.Chirp{}
	for _, c := range chirps {
		if c.ID == chirpID {
			chirp = c
			break
		}
	}

	if chirp.ID == 0 {
		respondWithError(w, http.StatusNotFound, "Chirp was not found")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
