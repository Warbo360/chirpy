package main

import (
	"context"
	"net/http"
	"time"

	"github.com/Warbo360/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerGetRefresh(w http.ResponseWriter, req *http.Request) {
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

	// Check if the token has been revoked if valid, it means the Time is valid and has been timestamped and thus
	// revoked, if valid is false then the time is not set and thus token has not been revoked

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token has been revoked", err)
		return
	}

	// Check if refresh token is expired or not

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "refresh token is expired", nil)
		return
	}

	// Query for the user based on the refresh token

	user, err := cfg.db.GetUserFromRefreshToken(context.Background(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting user from refresh token", err)
		return
	}

	// Create a new access token for the user

	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(cfg.accessTokenTimeoutSecounds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating access token", err)
		return
	}

	// Response to request is just the newly created access token

	respondWithJSON(w, http.StatusOK, returnVal{accessToken})
}
