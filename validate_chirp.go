package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type RequestBody struct {
	// the key will be the name of struct field unless you give it an explicit JSON tag
	Body string `json:"body"`
}

type ResponseData struct {
	// the key will be the name of struct field unless you give it an explicit JSON tag
	Body string `json:"cleaned_body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	// Create a variable to hold the JSON data.
	var body RequestBody

	// Decode the JSON data from the request body.
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	if len(body.Body) > 140 {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body.Body, badWords)

	responseBody := ResponseData{Body: cleaned}
	respondWithJSON(w, 200, responseBody)
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
