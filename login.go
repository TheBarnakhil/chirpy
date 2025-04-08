package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/TheBarnakhil/chirpy/internal/auth"
	"github.com/TheBarnakhil/chirpy/internal/database"
)

func (cfg *apiConfig) login(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "couldn't read request", err)
		return
	}

	rb := requestBody{}
	if err := json.Unmarshal(data, &rb); err != nil {
		respondWithError(rw, http.StatusInternalServerError, "couldn't unmarshal req body", err)
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), rb.Email)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if err := auth.CheckHashPassword(user.HashedPassword, rb.Password); err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)

	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error generating token", err)
		return
	}

	refToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error generating refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(
		req.Context(),
		database.CreateRefreshTokenParams{
			Token:     refToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		},
	)

	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error storing refresh token", err)
		return
	}

	respondWithJson(rw, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refToken,
	})
}
