package main

import "net/http"

func healthHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

	rw.WriteHeader(200)
	rw.Write([]byte("OK"))
}
