package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/TheBarnakhil/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading env from .env file: %v", err)
	}

	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}

	serveMux := http.NewServeMux()
	apiCfg := &apiConfig{fileServerHits: atomic.Int32{}, db: database.New(db), platform: os.Getenv("PLATFORM"), secret: os.Getenv("SECRET_KEY")}

	fileServer := http.FileServer(http.Dir(filepathRoot))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	serveMux.HandleFunc("GET /api/healthz", healthHandler)

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.hitsHander)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	serveMux.HandleFunc("POST /api/login", apiCfg.login)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)

	serveMux.HandleFunc("POST /api/users", apiCfg.createUser)
	serveMux.HandleFunc("PUT /api/users", apiCfg.updateUser)

	serveMux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	serveMux.HandleFunc("GET  /api/chirps", apiCfg.getAllChirps)
	serveMux.HandleFunc("GET  /api/chirps/{chirpID}", apiCfg.getChirpById)
	serveMux.HandleFunc("DELETE  /api/chirps/{chirpID}", apiCfg.delChirpById)

	serveMux.HandleFunc("POST /api/polka/webhooks", apiCfg.polkaWebhook)

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
