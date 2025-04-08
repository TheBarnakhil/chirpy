package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhook(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "couldn't read request", err)
		return
	}

	rb := requestBody{}
	if err := json.Unmarshal(data, &rb); err != nil {
		respondWithError(rw, http.StatusInternalServerError, "couldn't unmarshal req body", err)
		return
	}

	if rb.Event != "user.upgraded" {
		respondWithJson(rw, http.StatusNoContent, struct{}{})
		return
	}
	user, err := cfg.db.UpgradeUserToRed(req.Context(), rb.Data.UserID)
	if user.ID == uuid.Nil {
		respondWithError(rw, http.StatusNotFound, "no user found", err)
		return
	}
	if err != nil {
		respondWithError(rw, http.StatusNotFound, "couldn't upgrade to red", err)
		return
	}

	respondWithJson(rw, http.StatusNoContent, struct{}{})
}
