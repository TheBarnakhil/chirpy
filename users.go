package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/TheBarnakhil/chirpy/internal/auth"
	"github.com/TheBarnakhil/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUser(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
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

	hashed, err := auth.HashPassword(rb.Password)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "couldn't hash pass", err)
		return
	}

	user, err := cfg.db.CreateUser(
		req.Context(),
		database.CreateUserParams{
			Email:          rb.Email,
			HashedPassword: hashed,
		})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "error when creating user", err)
		return
	}

	respondWithJson(rw, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}

func (cfg *apiConfig) updateUser(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "couldn't get bearer token", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "couldn't validate token", err)
		return
	}

	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
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

	hashed, err := auth.HashPassword(rb.Password)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "couldn't hash pass", err)
		return
	}

	user, err := cfg.db.UpdateUserById(
		req.Context(),
		database.UpdateUserByIdParams{
			Email:          rb.Email,
			HashedPassword: hashed,
			ID:             id,
		})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "error when creating user", err)
		return
	}

	respondWithJson(rw, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
