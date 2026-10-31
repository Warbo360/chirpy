package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Warbo360/chirpy/internal/auth"
	"github.com/Warbo360/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "No email provided", nil)
		return
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}
	params.Password = hash
	dbusr, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: params.Password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failure creating user into database", err)
		return
	}
	usr := User{
		ID:        dbusr.ID,
		CreatedAt: dbusr.CreatedAt,
		UpdatedAt: dbusr.UpdatedAt,
		Email:     dbusr.Email,
	}
	respondWithJSON(w, http.StatusCreated, usr)
}
