package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hatimhas/chirtpy/internal/database"
)

func (c *apiConfig) handlerAddChirps(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := parameters{}

	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldnt Decode parameters", err)
	}

	if len(reqParams.Body) > 140 {
		respondWithErr(w, http.StatusBadRequest, "Chirp is too long", nil)
		return

	}

	cleanedreqBody := profaneCheck(reqParams.Body)

	chirp, err := c.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedreqBody,
		UserID: reqParams.UserId,
	})
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldnt create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	},
	)
}
