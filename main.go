package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(filepathRoot))
	serveMux.Handle("/app/", http.StripPrefix("/app", fileServer))
	serveMux.HandleFunc("/healthz", healthHandler)

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func healthHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

	rw.WriteHeader(200)
	rw.Write([]byte("OK"))
}
