package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (c *apiConfig) handlerGetChirpByID(w http.ResponseWriter, req *http.Request) {
	chirpReqId := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpReqId)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "error parsing chirpID to uuid", err)
		return
	}
	chirpDB, err := c.dbQueries.GetChirpById(req.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithErr(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithErr(w, http.StatusInternalServerError, "Failed to retrieve chirp by chirpID", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body:      chirpDB.Body,
		UserId:    chirpDB.UserID,
	})
}
