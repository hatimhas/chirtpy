package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	dbChirps, err := c.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldnt get all chirps", err)
	}
	chirps := []Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
