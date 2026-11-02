package main

import (
	"context"
	"net/http"

	"github.com/Warbo360/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, req *http.Request) {
	// Struct to return access token to requester

	type returnVal struct {
		Token string `json:"token"`
	}

	// Get refresh token string from the Authorization key in the header

	refreshTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error obtaining refresh token", err)
		return
	}

	// Query database for the retrived token string

	refreshToken, err := cfg.db.GetRefreshToken(context.Background(), refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting refresh token for database", err)
		return
	}

	// Query that updates revoked_at time

	err = cfg.db.UpdateRefreshTokenRevokeTime(context.Background(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating revoked_at time", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
