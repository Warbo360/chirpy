package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Warbo360/chirpy/internal/auth"
	"github.com/Warbo360/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, req *http.Request) {
	// parameters struct to decode as for query

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// New response type since User type does not cover JWT tokens here

	type returnVal struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Decode request body into above struct

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding login parameters", err)
		return
	}

	// Check for empty email and password

	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "no email provided", nil)
		return
	}

	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "no password provided", nil)
		return
	}

	// Get user by email from database

	dbUser, err := cfg.db.GetUser(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// If user hit in DB, check passed password against DB hashed password

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error verifying user password", err)
		return
	}

	// If password does not match respond with error

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	accessTokenNew, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Duration(cfg.accessTokenTimeoutSecounds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error during token generation", err)
		return
	}

	// Refresh Token Generation and insert into the database

	refreshTokenString := auth.MakeRefreshToken()
	refreshTokenNew, err := cfg.db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:  refreshTokenString,
		UserID: dbUser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token", err)
		return
	}

	// If password does match respond with User db fields minus password

	respondWithJSON(w, http.StatusOK, returnVal{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
		Token:        accessTokenNew,
		RefreshToken: refreshTokenNew.Token,
	})
}
