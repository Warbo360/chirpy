package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Warbo360/chirpy/internal/auth"
	"github.com/Warbo360/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateUserLogin(w http.ResponseWriter, req *http.Request) {
	// Parse the request header for the auth value which should be the access token

	accessTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error parsing access token", err)
		return
	}

	// Validate JWT access token and get user from token

	userID, err := auth.ValidateJWT(accessTokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error validating access token", err)
		return
	}

	// Define parameters struct to extract JSON body into

	type parameters struct {
		NewEmail          string `json:"email"`
		NewPassword       string `json:"password"`
		ID                uuid.UUID
		NewHashedPassword string
	}

	// Extrct JSON body into predefined parameters struct

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request body", err)
		return
	}

	// Pass user id to the parameters struct

	params.ID = userID
	params.NewHashedPassword, err = auth.HashPassword(params.NewPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating new hashed password", err)
		return
	}

	// Query to update user's email and password

	dbUser, err := cfg.db.UpdateUserEmailAndPassword(context.Background(), database.UpdateUserEmailAndPasswordParams{
		Email:          params.NewEmail,
		HashedPassword: params.NewHashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user's email and password in database", err)
		return
	}

	// Respond with succesfull update with User struct

	respondWithJSON(w, http.StatusOK, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	})
}
