package main

import (
	"net/http"

	"github.com/TheBarnakhil/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(rw http.ResponseWriter, req *http.Request) {
	type responseBody struct{}

	token, err := auth.GetBearerToken(req.Header)

	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error fetching bearer token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error an revoking refresh token", err)
		return
	}

	respondWithJson(rw, http.StatusNoContent, responseBody{})
}
