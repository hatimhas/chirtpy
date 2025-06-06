package main

import (
	"net/http"
)

func (c *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
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
