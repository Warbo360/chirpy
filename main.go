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
	fileserverHits             atomic.Int32
	db                         *database.Queries
	platform                   string
	secret                     string
	accessTokenTimeoutSecounds int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	scrt := os.Getenv("SECRET")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	pltfrm := os.Getenv("PLATFORM")
	if pltfrm == "" {
		log.Fatal("PLATFORM must be set")
	}
	if scrt == "" {
		log.Fatal("SECRET must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failure connecting to postgres database")
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits:             atomic.Int32{},
		db:                         dbQueries,
		platform:                   pltfrm,
		secret:                     scrt,
		accessTokenTimeoutSecounds: 3600,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsLoggerHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.hitsResetHandler)
	mux.HandleFunc("POST /api/users/", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerGetRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUserLogin)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirpByID)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %v on port: %v\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
