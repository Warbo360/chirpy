package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, resp *http.Request) {
	// Get Chirp ID to search for in request URL

	chirpID := resp.PathValue("id")
	if chirpID == "" {
		respondWithError(w, http.StatusNotFound, "no ID provided in request", nil)
		return
	}

	// Parse Chirp ID from string to UUID

	parsedChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error parsing Chirp ID from request", err)
		return
	}

	// Query database for chirp by ID

	dbChirp, err := cfg.db.GetChirp(context.Background(), parsedChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error searching database for Chirp ID", err)
		return
	}

	// Map database Chirp to Go Chirp struct with JSON struct tags

	goChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, goChirp)
}
