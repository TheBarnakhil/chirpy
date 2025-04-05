package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits.Add(1) // Increment the counter
		next.ServeHTTP(rw, req)   // Call the next handler
	})
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	apiCfg := &apiConfig{fileServerHits: atomic.Int32{}}

	fileServer := http.FileServer(http.Dir(filepathRoot))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	serveMux.HandleFunc("/healthz", healthHandler)
	serveMux.HandleFunc("/metrics", apiCfg.hitsHander)
	serveMux.HandleFunc("/reset", apiCfg.resetHandler)

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) hitsHander(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

	rw.WriteHeader(200)
	rw.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())))
}

func (cfg *apiConfig) resetHandler(rw http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Store(0)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

	rw.WriteHeader(200)
	rw.Write([]byte("Reset the hits"))
}
