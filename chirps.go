package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/admiralyeoj/chirpy/internal/database"
)

const dbPath string = "database.json"

type RequestBody struct {
	// the key will be the name of struct field unless you give it an explicit JSON tag
	Body string `json:"body"`
}

func handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// Create a variable to hold the JSON data.
	var body RequestBody
	db, err := database.NewDB(dbPath)

	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	// Decode the JSON data from the request body.
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	if len(body.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body.Body, badWords)

	chirp, err := db.CreateChirp(cleaned)

	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 201, chirp)
}

func handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	// Create a variable to hold the JSON data.
	db, err := database.NewDB(dbPath)

	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	chirps, err := db.GetChirps()

	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	respondWithJSON(w, 200, chirps)
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
