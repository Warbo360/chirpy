package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateUserToIsRed(w http.ResponseWriter, req *http.Request) {
	// Define parameters struct we expect in the request body

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	// Decode request body into a params struct

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request JSON", err)
		return
	}

	// If event passed to endpoint is anything other than "user.upgraded" respond with not content, we do not care for
	// any event past as such.

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	parsedUserID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing request for user ID", err)
		return
	}

	err = cfg.db.UpdateUserIsRedToTrue(context.Background(), parsedUserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error upgrading user in database", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
