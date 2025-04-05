package main

import (
	"log"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}
	fileServer := http.FileServer(http.Dir('.'))
	serveMux.Handle("/", fileServer)
	log.Fatal(server.ListenAndServe())
}
