package main

import (
	"github.com/admiralyeoj/chirpy/internal/database"
)

const dbPath string = "database.json"

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}
