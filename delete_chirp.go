package main

import (
	"context"
	"net/http"

	"github.com/Warbo360/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, req *http.Request) {
	// Extract chirp ID from path value

	chirpID := req.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusBadRequest, "no chirp ID provided in request pathway", nil)
		return
	}
	parsedChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing Chirp ID", err)
		return
	}

	// Parse request header for authentication value which should be the access JWT token

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error parsing access token", err)
		return
	}

	// Validate the parsed token and get the user ID

	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error validating access token", err)
		return
	}

	// Query to get chirp by chirp ID

	dbChirp, err := cfg.db.GetChirp(context.Background(), parsedChirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting Chirp from database", err)
		return
	}

	// Check if dbChirp.UserID matches UserID hit from JWT validation if not respond with error, if it does query to
	// delete chirp entry in chirps table in database

	if userID != dbChirp.UserID {
		respondWithError(w, http.StatusForbidden, "user not authorized", err)
		return
	}
	err = cfg.db.DeleteChirpByID(context.Background(), dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting Chirp from database", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
