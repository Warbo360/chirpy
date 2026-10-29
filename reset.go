package main

import (
	"context"
	"net/http"
)

// Method to reset server hits

func (cfg *apiConfig) hitsResetHandler(w http.ResponseWriter, req *http.Request) {
	err := cfg.db.ResetUsers(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failure to reset \"users\" table", err)
		return
	}
	respondWithJSON(w, http.StatusOK, "Users table reset successfully")
}
