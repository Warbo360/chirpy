package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
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
	dbusr, err := cfg.db.CreateUser(context.Background(), params.Email)
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
