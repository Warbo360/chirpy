package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Warbo360/chirpy/internal/database"
	"github.com/google/uuid"
)

const maxChirpLen = 140

var profanes = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	// Parameters the handler takes in from req JSON

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// JSON unmarhsaling into parameters

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters", err)
		return
	}

	// Validates the parameters body for length and profanity

	cleanedBody, err := helperChirpValidate(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failure in validating chirp body", err)
		return
	}
	params.Body = cleanedBody

	// DB query to create chirp into the chirps DB with the given parameters

	dbchirp, err := cfg.db.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failure in creating chirp to database", err)
		return
	}

	// Internal DB chirp fields to Go representation of that with struct tags to go back to JSON in the response

	chirp := Chirp{
		ID:        dbchirp.ID,
		CreatedAt: dbchirp.CreatedAt,
		UpdatedAt: dbchirp.UpdatedAt,
		Body:      dbchirp.Body,
		UserID:    dbchirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

// Helper function that validates chirps for length and profanity

func helperChirpValidate(body string) (string, error) {
	if len(body) > maxChirpLen {
		return "", errors.New("Chirp is too long")
	}
	cleaned := helperChirpSanitize(body)
	return cleaned, nil
}

// Helper function that cleans chirps from list profanity

func helperChirpSanitize(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := profanes[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
