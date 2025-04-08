package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits.Add(1) // Increment the counter
		next.ServeHTTP(rw, req)   // Call the next handler
	})
}

func (cfg *apiConfig) hitsHander(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	rw.WriteHeader(http.StatusOK)
	template := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	rw.Write([]byte(fmt.Sprintf(template, cfg.fileServerHits.Load())))
}
