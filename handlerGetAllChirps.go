package main

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/hatimhas/chirtpy/internal/database"
)

func (c *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	authorIDstr := req.URL.Query().Get("author_id")
	sortstr := req.URL.Query().Get("sort")

	dbChirps := []database.Chirp{}
	var err error

	if authorIDstr == "" {
		dbChirps, err = c.dbQueries.GetAllChirps(req.Context())
		if err != nil {
			respondWithErr(w, http.StatusInternalServerError, "Failed Retreiving All Chirps", err)
			return
		}

	} else {

		authorID, err := uuid.Parse(authorIDstr)
		if err != nil {
			respondWithErr(w, http.StatusBadRequest, fmt.Sprintf("author_id %s is not a valid uuid", authorIDstr), err)
			return
		}

		dbChirps, err = c.dbQueries.GetChirpsByAuthorId(req.Context(), authorID)
		if err != nil {
			respondWithErr(w, http.StatusInternalServerError, "Couldnt get chirps by authorID", err)
			return
		}

	}

	chirps := make([]Chirp, len(dbChirps))
	for i, chirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		}
	}

	if sortstr == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
