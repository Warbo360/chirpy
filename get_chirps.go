package main

import (
	"context"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	// Check for an author_id in the request URL, if author ID exists parse as a user ID and get all chirps by user ID

	authorID := req.URL.Query().Get("author_id")
	sortOrder := req.URL.Query().Get("sort")
	if sortOrder == "" {
		sortOrder = "asc"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		respondWithError(w, http.StatusBadRequest, "not supported sorting in query parameter", nil)
		return
	}
	if authorID != "" {
		parsedAuthorID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing author ID", err)
			return
		}
		dbChirps, err := cfg.db.GetChirpsByUser(context.Background(), parsedAuthorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error retrieving Chirps by passed author ID", err)
			return
		}
		var chirps []Chirp
		for _, chirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}
		if sortOrder == "asc" {
			sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
		} else {
			sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
		}
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	// Query that returns all chirps from database

	dbChirps, err := cfg.db.GetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failure retrieving chirps from database", err)
		return
	}

	// Loop through all chirps and append to a slice of chirps to include in the JSON payload

	var chirps []Chirp
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	if sortOrder == "asc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
