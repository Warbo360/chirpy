package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

var profanes = []string{"kerfuffle", "sharbert", "fornax"}

type parameters struct {
	Body string `json:"body"`
}

type returnVals struct {
	Valid bool `json:"valid"`
}

func handlerChirpValidate(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	const maxChirpLen = 140
	if len(params.Body) > maxChirpLen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, returnVals{Valid: true})
}

func chirpSanitize(chirp string) string {
	sChirp := strings.Split(chirp, " ")
	for i, sub := range sChirp {
		if slices.Contains(profanes, strings.ToLower(sub)) {
			sChirp[i] = "****"
		}
	}
	chirp = strings.Join(sChirp, " ")
	return chirp
}
