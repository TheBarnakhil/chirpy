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

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirp(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		Chirp
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "couldn't read request", err)
		return
	}

	reqBody := requestBody{}
	if err := json.Unmarshal(data, &reqBody); err != nil {
		respondWithError(rw, http.StatusInternalServerError, "couldn't unmarshal req body", err)
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "couldn't get bearer token", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "couldn't validate token", err)
		return
	}

	if len(reqBody.Body) <= 140 {
		cleanedBody, err := filterProfaneWords(reqBody.Body)
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, "error cleaning body", err)
		}
		chirp, err := cfg.db.CreateChirp(
			req.Context(),
			database.CreateChirpParams{
				Body:   cleanedBody,
				UserID: id,
			})
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, "error creating chirp", err)
		}
		respondWithJson(rw, http.StatusCreated,
			responseBody{Chirp: Chirp{
				ID:        chirp.ID,
				Body:      chirp.Body,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				UserID:    chirp.UserID,
			}})
	} else {
		respondWithError(rw, 400, "Chirp is too long", err)
	}
}

func (cfg *apiConfig) getAllChirps(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	dbChirps, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error retrieving chirps from db", err)
		return
	}

	chirps := []Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID:    chirp.UserID,
			Body:      chirp.Body,
		})
	}

	respondWithJson(rw, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpById(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type responseBody struct {
		Chirp
	}

	val, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error parsing chirpId", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(req.Context(), val)
	if err != nil {
		respondWithError(rw, http.StatusNotFound, "Error retrieving chirp from db", err)
		return
	}

	respondWithJson(rw, http.StatusOK, responseBody{
		Chirp: Chirp{
			ID:        chirp.ID,
			Body:      chirp.Body,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID:    chirp.UserID,
		}})
}

func (cfg *apiConfig) delChirpById(rw http.ResponseWriter, req *http.Request) {
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

	type responseBody struct{}

	val, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error parsing chirpId", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(req.Context(), val)
	if err != nil {
		respondWithError(rw, http.StatusNotFound, "Error retrieving chirp from db", err)
		return
	}

	if chirp.UserID != id {
		respondWithError(rw, http.StatusForbidden, "Error not authorized", err)
		return
	}

	err = cfg.db.DeleteChirpById(req.Context(), val)
	if err != nil {
		respondWithError(rw, http.StatusNotFound, "Error retrieving chirp from db", err)
		return
	}

	respondWithJson(rw, http.StatusNoContent, responseBody{})
}
