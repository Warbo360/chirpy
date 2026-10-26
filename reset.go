package main

import (
	"net/http"
)

// Method to reset server hits

func (cfg *apiConfig) hitsResetHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server hits reset to 0\n"))
}
