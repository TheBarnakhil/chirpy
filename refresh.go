package main

import (
	"net/http"
	"time"

	"github.com/TheBarnakhil/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(rw http.ResponseWriter, req *http.Request) {
	type responseBody struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(req.Header)

	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error fetching bearer token", err)
		return
	}

	refTok, err := cfg.db.GetRefreshToken(req.Context(), token)
	if err != nil || time.Now().Local().After(refTok.ExpiresAt) || refTok.RevokedAt.Valid {
		respondWithError(rw, http.StatusUnauthorized, "This is an invalid refresh token", err)
		return
	}

	accessTok, err := auth.MakeJWT(refTok.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error fetching bearer token", err)
		return
	}

	respondWithJson(rw, http.StatusOK, responseBody{Token: accessTok})
}
