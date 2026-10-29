package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Warbo360/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Struct for holding stateful server data

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	pltfrm := os.Getenv("PLATFORM")
	if pltfrm == "" {
		log.Fatal("PLATFORM must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failure connecting to postgres database")
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       pltfrm,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsLoggerHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.hitsResetHandler)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpValidate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %v on port: %v\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
