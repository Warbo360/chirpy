package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

// Struct for holding stateful server data

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	var apiCfg apiConfig
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsLoggerHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.hitsResetHandler)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %v on port: %v\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
