package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/admiralyeoj/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// Create a variable to hold the JSON data.
	var body RequestBody

	// Decode the JSON data from the request body.
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(body.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body.Body, badWords)

	chirp, err := cfg.DB.CreateChirp(cleaned)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, database.Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
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

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
