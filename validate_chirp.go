package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
		respondWithError(w, 500, "Error decoding parameters", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	respBody := returnVals{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		respondWithError(w, 500, "Error encoding response body", err)
		return
	}
	respondWithJSON(w, 200, dat)
	return
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err == nil {
		log.Printf("%s", msg)
	} else {
		log.Printf("%s: %s", msg, err)
	}
	w.WriteHeader(code)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data, ok := payload.([]byte); ok {
		w.Write(data)
	}
}
