package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	// Single return stuct for

	chirpsdb, err := cfg.db.GetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failure retrieving chirps from database", err)
		return
	}

	var chirps []Chirp
	for _, chirp := range chirpsdb {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
