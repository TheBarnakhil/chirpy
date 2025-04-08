package main

import (
	"errors"
	"net/http"
)

func (cfg *apiConfig) resetHandler(rw http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(rw, http.StatusForbidden, "You have stumbled upon a forbidden zone!", errors.New("only available in dev"))
		return
	}

	cfg.fileServerHits.Store(0)
	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error deleting all users", err)
		return
	}
	err = cfg.db.DeleteChirps(req.Context())
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error deleting all chirps", err)
		return
	}
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
