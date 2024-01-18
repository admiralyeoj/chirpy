package main

import (
	"github.com/admiralyeoj/chirpy/internal/database"
)

const dbPath string = "database.json"

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

type RequestBody struct {
	// the key will be the name of struct field unless you give it an explicit JSON tag
	Body  string `json:"body"`
	Email string `json:"email"`
}
