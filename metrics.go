package main

import (
	"fmt"
	"net/http"
)

// Middleware to increment server hits on being handled

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

// Method to print server hits

func (cfg *apiConfig) hitsLoggerHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	var buf []byte

	// HTML to load on being handled to admin metrics.... This is gross

	buf = fmt.Appendf(buf, `
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		<html>`,
		cfg.fileserverHits.Load())

	w.Write(buf)
}
