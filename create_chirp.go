package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Warbo360/chirpy/internal/auth"
	"github.com/Warbo360/chirpy/internal/database"
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
		Body string `json:"body"`
	}

	// JSON unmarhsaling into parameters

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters", err)
		return
	}

	// Getting JWT from http.Header in the request

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting token from header", err)
		return
	}
	user, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error verifying user token", err)
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
		UserID: user,
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
